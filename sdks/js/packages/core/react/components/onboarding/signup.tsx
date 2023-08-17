import { Flex, Link, Text } from '@raystack/apsara';
import React, { ComponentPropsWithRef, useCallback } from 'react';
import { useFrontier } from '~/react/contexts/FrontierContext';
import { Container } from '../Container';
import { Header } from '../Header';
import { MagicLink } from './magiclink';
import { OIDCButton } from './oidc';

const styles = {
  titleContainer: {
    fontWeight: '400'
  }
};

type SignUpProps = ComponentPropsWithRef<typeof Container> & {
  logo?: React.ReactNode;
  title?: string;
  excludes?: string[];
};
export const SignUp = ({
  logo,
  title = 'Create your account',
  excludes = [],
  ...props
}: SignUpProps) => {
  const { config } = useFrontier();
  const { client, strategies = [] } = useFrontier();

  const clickHandler = useCallback(
    async (name?: string) => {
      if (!name) return;
      if (!client) return;

      const {
        data: { endpoint = '' }
      } = await client.frontierServiceAuthenticate(name);

      window.location.href = endpoint;
    },
    [strategies]
  );

  const mailotp = strategies.find(s => s.name === 'mailotp');
  const filteredOIDC = strategies
    .filter(s => s.name !== 'mailotp')
    .filter(s => !excludes.includes(s.name ?? ''));

  return (
    <Container {...props}>
      <Header logo={logo} title={title} />
      <Flex direction="column" style={{ width: '100%', gap: '8px' }}>
        {filteredOIDC.map((s, index) => {
          return (
            <OIDCButton
              key={index}
              onClick={() => clickHandler(s.name)}
              provider={s.name || ''}
            ></OIDCButton>
          );
        })}

        {mailotp && <MagicLink />}
      </Flex>
      <div style={styles.titleContainer}>
        <Text size={2}>
          Already have an account?{' '}
          <Link
            href={config.redirectLogin}
            style={{ color: 'var(--foreground-accent)' }}
          >
            Login
          </Link>
        </Text>
      </div>
    </Container>
  );
};