/*
Copyright 2019 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import 'xterm/css/xterm.css';
import { IDisposable, ITheme, Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import { debounce } from 'shared/utils/highbar';

import { IPtyProcess } from 'teleterm/sharedProcess/ptyHost';
import Logger from 'teleterm/logger';

const WINDOW_RESIZE_DEBOUNCE_DELAY = 200;

type Options = {
  el: HTMLElement;
  fontSize: number;
  theme: ITheme;
};

export default class TtyTerminal {
  public term: Terminal;
  private el: HTMLElement;
  private fitAddon = new FitAddon();
  private resizeHandler: IDisposable;
  private debouncedResize: () => void;
  private logger = new Logger('lib/term/terminal');
  private removePtyProcessOnDataListener: () => void;

  constructor(private ptyProcess: IPtyProcess, private options: Options) {
    this.el = options.el;
    this.term = null;

    this.debouncedResize = debounce(
      this.requestResize.bind(this),
      WINDOW_RESIZE_DEBOUNCE_DELAY
    );
  }

  open(): void {
    this.term = new Terminal({
      cursorBlink: false,
      /**
       * `fontFamily` can be provided by the user and is unsanitized. This means that it cannot be directly used in CSS,
       * as it may inject malicious CSS code.
       * To sanitize the value, we set it as a style on the HTML element and then read it from it.
       * Read more https://frontarm.com/james-k-nelson/how-can-i-use-css-in-js-securely/.
       */
      fontFamily: this.el.style.fontFamily,
      fontSize: this.options.fontSize,
      scrollback: 5000,
      minimumContrastRatio: 4.5, // minimum for WCAG AA compliance
      theme: this.options.theme,
      windowOptions: {
        setWinSizeChars: true,
      },
    });

    this.term.loadAddon(this.fitAddon);

    this.registerResizeHandler();

    this.term.open(this.el);

    this.fitAddon.fit();

    this.term.onData(data => {
      this.ptyProcess.write(data);
    });

    this.term.onResize(size => {
      this.ptyProcess.resize(size.cols, size.rows);
    });

    this.removePtyProcessOnDataListener = this.ptyProcess.onData(data =>
      this.handleData(data)
    );

    // TODO(ravicious): Don't call start if the process was already started.
    // This is what is causing the terminal to visually repeat the input on hot reload.
    // The shared process version of PtyProcess knows whether it was started or not (the status
    // field), so it's a matter of exposing this field through gRPC and reading it here.
    this.ptyProcess.start(this.term.cols, this.term.rows);

    window.addEventListener('resize', this.debouncedResize);
  }

  focus() {
    this.term.focus();
  }

  requestResize(): void {
    const visible = !!this.el.clientWidth && !!this.el.clientHeight;
    if (!visible) {
      this.logger.info(`unable to resize terminal (container might be hidden)`);
      return;
    }
    this.fitAddon.fit();
  }

  destroy(): void {
    this.removePtyProcessOnDataListener?.();
    this.term?.dispose();
    this.fitAddon.dispose();
    this.resizeHandler?.dispose();
    this.el.innerHTML = null;

    window.removeEventListener('resize', this.debouncedResize);
  }

  private registerResizeHandler(): void {
    let prevCols: number, prevRows: number;
    this.resizeHandler = this.term.parser.registerCsiHandler(
      { final: 't' },
      params => {
        const [ps] = params;
        if (ps === 8) {
          // Ps = 8 - resizes the text area to given height and width in characters.
          const rows = params[1] as number;
          const cols = params[2] as number;
          if (prevRows !== rows || prevCols !== cols) {
            prevRows = rows;
            prevCols = cols;
            this.term.resize(cols, rows);
          }
          return true; // sequence has been handled
        }
        return false;
      }
    );
  }

  private handleData(data: string): void {
    try {
      this.term.write(data);
    } catch (err) {
      this.logger.error('xterm.write', data, err);
      // recover xtermjs by resetting it
      this.term.reset();
    }
  }
}
