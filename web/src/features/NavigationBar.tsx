import { IconBrandGithubFilled, IconHeartFilled, IconInfoCircle } from '@tabler/icons-react';
import { Link } from 'react-router-dom';
import { ActionIcon, Anchor, AppShell, Flex, Group, Image, Text, Tooltip } from '@mantine/core';
import logo from '../assets/logo.svg';
import { Footer } from './Footer';
import { LanguageSelector } from './language/LanguageSelector';

export const NavigationBar = ({ children }: { children: React.ReactNode }) => {
  return (
    <>
      <AppShell header={{ height: 60 }} padding="md">
        <AppShell.Header>
          <Group h="100%" px="md" justify="space-between">
            <Group>
              <Anchor component={Link} to="/" size="lg" underline="never" c="dark">
                <Flex direction="row" gap={8}>
                  <Image src={logo} h={30} w="auto" alt="Logo" />
                  <Text fw={700} size="lg" visibleFrom="xs">
                    Kubernetes i18n Tracker
                  </Text>
                </Flex>
              </Anchor>
            </Group>
            <Group gap={8}>
              {/* <ActionIcon
                component="a"
                href="/information"
                target="_blank"
                variant="default"
                radius="md"
                size="lg"
                aria-label="GitHub"
              >
                <Tooltip label="Information" position="bottom" withArrow offset={10}>
                  <IconInfoCircle size={20} />
                </Tooltip>
              </ActionIcon> */}
              <ActionIcon
                component="a"
                href="https://github.com/kfess/kubernetes-i18n-tracker"
                target="_blank"
                variant="default"
                radius="md"
                size="lg"
                aria-label="GitHub"
              >
                <Tooltip label="Source Code" position="bottom" withArrow offset={10}>
                  <IconBrandGithubFilled size={20} />
                </Tooltip>
              </ActionIcon>
              <ActionIcon
                component="a"
                href="https://github.com/kubernetes/website"
                target="_blank"
                variant="default"
                radius="md"
                size="lg"
                c="pink"
                aria-label="contribute-to-the-translation"
              >
                <Tooltip
                  label="Contribute to the translation"
                  position="bottom"
                  withArrow
                  offset={10}
                >
                  <IconHeartFilled size={20} />
                </Tooltip>
              </ActionIcon>
              <LanguageSelector />
            </Group>
          </Group>
        </AppShell.Header>
        <AppShell.Main>{children}</AppShell.Main>
      </AppShell>
      <Footer />
    </>
  );
};
