import { AppShell, Header, Navbar, Text, MediaQuery, Burger, useMantineTheme } from '@mantine/core';
import { useState } from 'react';
import { Horse, HeartRateMonitor, Baby } from 'tabler-icons-react';

export function Layout({ children }) {
    const theme = useMantineTheme();
    const [opened, setOpened] = useState(false);

    return (
        <AppShell
            styles={{
                main: {
                    background: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0],
                },
            }}
            navbarOffsetBreakpoint="sm"
            navbar={
                <Navbar
                    p="md"
                    hiddenBreakpoint="sm"
                    hidden={!opened}
                    width={{ sm: 200, lg: 300 }}
                >
                    <Navbar.Section>
                        <Text
                            component="a"
                            href="/"
                            size="sm"
                            weight={500}
                            sx={{ display: 'block', marginBottom: theme.spacing.xs }}
                        >
                            <Horse style={{ marginRight: 8 }} />
                            Horses
                        </Text>
                        <Text
                            component="a"
                            href="/health"
                            size="sm"
                            weight={500}
                            sx={{ display: 'block', marginBottom: theme.spacing.xs }}
                        >
                            <HeartRateMonitor style={{ marginRight: 8 }} />
                            Health Records
                        </Text>
                        <Text
                            component="a"
                            href="/pregnancy"
                            size="sm"
                            weight={500}
                            sx={{ display: 'block', marginBottom: theme.spacing.xs }}
                        >
                            <Baby style={{ marginRight: 8 }} />
                            Pregnancy Events
                        </Text>
                    </Navbar.Section>
                </Navbar>
            }
            header={
                <Header height={{ base: 50, md: 70 }} p="md">
                    <div style={{ display: 'flex', alignItems: 'center', height: '100%' }}>
                        <MediaQuery largerThan="sm" styles={{ display: 'none' }}>
                            <Burger
                                opened={opened}
                                onClick={() => setOpened((o) => !o)}
                                size="sm"
                                color={theme.colors.gray[6]}
                                mr="xl"
                            />
                        </MediaQuery>

                        <Text>Horse Tracking App</Text>
                    </div>
                </Header>
            }
        >
            {children}
        </AppShell>
    );
}
