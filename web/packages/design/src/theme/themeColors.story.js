/*
Copyright 2020 Gravitational, Inc.

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

import React from 'react';
import { useTheme } from 'styled-components';

import { Flex, Box, Text, Link } from '..';

export default {
  title: 'Design/Theme/Colors',
};

export const Colors = () => <ColorsComponent />;

const ColorsComponent = () => {
  const theme = useTheme();

  return (
    <Flex flexDirection="column" p="4">
      <Flex flexDirection="column">
        <Text mb="2" typography="h4">
          Levels
        </Text>
        <Text mb="2">
          Levels are used to reflect the perceived depth of elements in the UI.
          The further back an element is, the more "sunken" it is, and the more
          forwards it is, the more "elevated" it is (think CSS z-index). <br />A
          "sunken" colour would be used to represent something like the
          background of the app. While "surface" would be the colour of the
          primary surface where most content is located (such as tables). Any
          colours more "elevated" than that would be used for things such as
          popovers, menus, and dialogs. <br />
          You can read more on this concept{' '}
          <Link
            target="_blank"
            href="https://m3.material.io/styles/elevation/applying-elevation"
          >
            here.
          </Link>
        </Text>
        <ColorsBox mb="4" colors={theme.colors.levels} themeType="levels" />
        <Text mb="1" typography="h4">
          Spot Backgrounds
        </Text>
        <Text mb="2">
          Spot backgrounds are used as highlights or accents. They are not solid
          colours, instead, they are a slightly transparent mask that highlights
          the colour behind it. This makes them quite versatile and they can be
          used to accentuate or highlight components on any background.
          <br />
          They may be used for things such as indicating a hover or active state
          for an item in a menu, or for minor accents such as dividers. <br />
        </Text>
        <Flex>
          <SingleColorBox
            mb="4"
            mr={0}
            color={theme.colors.spotBackground[0]}
            themeType="levels"
            path="theme.colors.spotBackground[0]"
          />
          <SingleColorBox
            mb="4"
            mr={0}
            color={theme.colors.spotBackground[1]}
            themeType="levels"
            path="theme.colors.spotBackground[1]"
          />
          <SingleColorBox
            mb="4"
            mr={0}
            color={theme.colors.spotBackground[2]}
            themeType="levels"
            path="theme.colors.spotBackground[2]"
          />
        </Flex>
        <Text mb="2" typography="h4">
          Brand
        </Text>
        <SingleColorBox
          mb="4"
          path="theme.colors.brand"
          color={theme.colors.brand}
        />
        <Text mb="2" typography="h4">
          Shadows
        </Text>
        <Flex>
          <Box
            mb={4}
            mr={6}
            css={`
              height: 200px;
              width: 200px;
              border-radius: 8px;
              display: flex;
              justify-content: center;
              align-items: center;
              text-align: center;
              background: ${theme.colors.levels.surface};
              box-shadow: ${theme.boxShadow[0]};
            `}
          >
            <Text>theme.boxShadow[0]</Text>
          </Box>
          <Box
            mb={4}
            mr={6}
            css={`
              height: 200px;
              width: 200px;
              border-radius: 8px;
              display: flex;
              justify-content: center;
              align-items: center;
              text-align: center;
              background: ${theme.colors.levels.surface};
              box-shadow: ${theme.boxShadow[1]};
            `}
          >
            <Text>theme.boxShadow[1]</Text>
          </Box>
          <Box
            mb={4}
            mr={6}
            css={`
              height: 200px;
              width: 200px;
              border-radius: 8px;
              display: flex;
              justify-content: center;
              align-items: center;
              text-align: center;
              background: ${theme.colors.levels.surface};
              box-shadow: ${theme.boxShadow[2]};
            `}
          >
            <Text>theme.boxShadow[2]</Text>
          </Box>
        </Flex>
        <Text mb="2" typography="h4">
          Text Colors
        </Text>
        <Flex width="fit-content" flexDirection="row" mb={4}>
          <Flex
            flexDirection="column"
            border={`1px solid ${theme.colors.text.muted}`}
            bg={theme.colors.levels.surface}
            py={3}
            px={3}
            mr={3}
          >
            <Text>theme.colors.text.main</Text>
            <Text typography="h1" color={theme.colors.text.main}>
              Primary
            </Text>
            <Text typography="h2" color={theme.colors.text.main}>
              Primary
            </Text>
            <Text typography="h3" color={theme.colors.text.main}>
              Primary
            </Text>
            <Text typography="h4" color={theme.colors.text.main}>
              Primary
            </Text>
            <Text typography="paragraph" color={theme.colors.text.main}>
              Primary
            </Text>
          </Flex>
          <Flex
            flexDirection="column"
            border={`1px solid ${theme.colors.text.muted}`}
            bg={theme.colors.levels.surface}
            py={3}
            px={3}
            mr={3}
          >
            <Text>theme.colors.text.slightlyMuted</Text>
            <Text typography="h1" color={theme.colors.text.slightlyMuted}>
              Secondary
            </Text>
            <Text typography="h2" color={theme.colors.text.slightlyMuted}>
              Secondary
            </Text>
            <Text typography="h3" color={theme.colors.text.slightlyMuted}>
              Secondary
            </Text>
            <Text typography="h4" color={theme.colors.text.slightlyMuted}>
              Secondary
            </Text>
            <Text
              typography="paragraph"
              color={theme.colors.text.slightlyMuted}
            >
              Secondary
            </Text>
          </Flex>
          <Flex
            flexDirection="column"
            border={`1px solid ${theme.colors.text.muted}`}
            bg={theme.colors.levels.surface}
            py={3}
            px={3}
            mr={3}
          >
            <Text>theme.colors.text.muted</Text>
            <Text typography="h1" color={theme.colors.text.muted}>
              Placeholder
            </Text>
            <Text typography="h2" color={theme.colors.text.muted}>
              Placeholder
            </Text>
            <Text typography="h3" color={theme.colors.text.muted}>
              Placeholder
            </Text>
            <Text typography="h4" color={theme.colors.text.muted}>
              Placeholder
            </Text>
            <Text typography="paragraph" color={theme.colors.text.muted}>
              Placeholder
            </Text>
          </Flex>
          <Flex
            flexDirection="column"
            border={`1px solid ${theme.colors.text.muted}`}
            bg={theme.colors.levels.surface}
            py={3}
            px={3}
            mr={3}
          >
            <Text>theme.colors.text.disabled</Text>
            <Text typography="h1" color={theme.colors.text.disabled}>
              Disabled
            </Text>
            <Text typography="h2" color={theme.colors.text.disabled}>
              Disabled
            </Text>
            <Text typography="h3" color={theme.colors.text.disabled}>
              Disabled
            </Text>
            <Text typography="h4" color={theme.colors.text.disabled}>
              Disabled
            </Text>
            <Text typography="paragraph" color={theme.colors.text.disabled}>
              Disabled
            </Text>
          </Flex>
          <Flex
            flexDirection="column"
            border={`1px solid ${theme.colors.text.muted}`}
            bg={theme.colors.text.main}
            py={3}
            px={3}
            mr={3}
          >
            <Text color={theme.colors.text.primaryInverse}>
              theme.colors.text.primaryInverse
            </Text>
            <Text typography="h1" color={theme.colors.text.primaryInverse}>
              Primary Inverse
            </Text>
            <Text typography="h2" color={theme.colors.text.primaryInverse}>
              Primary Inverse
            </Text>
            <Text typography="h3" color={theme.colors.text.primaryInverse}>
              Primary Inverse
            </Text>
            <Text typography="h4" color={theme.colors.text.primaryInverse}>
              Primary Inverse
            </Text>
            <Text
              typography="paragraph"
              color={theme.colors.text.primaryInverse}
            >
              Primary Inverse
            </Text>
          </Flex>
        </Flex>
        <Text mb="2" typography="h3">
          Interactive Colors
        </Text>
        <Text mb="2">
          Interactive colors are used for hover states, indicators, etc. An
          example of this in use currently would be unified resource cards in
          the Pinned and Pinned(Hovered) states.
        </Text>
        <Flex gap={4}>
          <SingleColorBox
            mb="2"
            path="theme.colors.interactive.tonal.primary[0]"
            color={theme.colors.interactive.tonal.primary[0]}
          />
          <SingleColorBox
            mb="2"
            path="theme.colors.interactive.tonal.primary[1]"
            color={theme.colors.interactive.tonal.primary[1]}
          />
          <SingleColorBox
            mb="2"
            path="theme.colors.interactive.tonal.primary[2]"
            color={theme.colors.interactive.tonal.primary[2]}
          />
        </Flex>
      </Flex>
    </Flex>
  );
};

function ColorsBox({ colors, themeType = null, ...styles }) {
  const list = Object.keys(colors).map(key => {
    const fullPath = themeType
      ? `theme.colors.${themeType}.${key}`
      : `theme.colors.${key}`;

    return (
      <Flex flexWrap="wrap" key={key} width="260px" mb={3}>
        <Box
          css={`
            color: ${props => props.theme.colors.text.slightlyMuted};
          `}
        >
          {fullPath}
        </Box>
        <Box
          width="100%"
          height="50px"
          p={3}
          mr={3}
          css={`
            background: ${colors[key]};
            border: 1px solid ${props => props.theme.colors.primaryInverse};
          `}
        />
      </Flex>
    );
  });

  return (
    <Flex flexWrap="wrap" {...styles}>
      {list}
    </Flex>
  );
}

function SingleColorBox({ color, path, ...styles }) {
  return (
    <Flex flexWrap="wrap" key={path} width="260px" mb={3}>
      <Box
        css={`
          color: ${props => props.theme.colors.text.slightlyMuted};
        `}
      >
        {path}
      </Box>
      <Box width="100%" height="50px" p={3} mr={3} bg={color} {...styles} />
    </Flex>
  );
}
