import '@rainbow-me/rainbowkit/styles.css';
import { getDefaultConfig, RainbowKitProvider, ConnectButton } from '@rainbow-me/rainbowkit';
import { WagmiProvider } from 'wagmi';
import { mainnet } from 'wagmi/chains';
import { QueryClientProvider, QueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { Button } from 'react-bootstrap';

const config = getDefaultConfig({
  appName: 'Apache Answer',
  projectId: 'xxx',
  chains: [mainnet],  // There's no on-chain operations, so only ETH mainnet is enough.
});

const queryClient = new QueryClient();

function WalletAuthorizer() {
  const { t } = useTranslation('plugin', {
    keyPrefix: 'connector_wallet_route.frontend',
  });

  return (
    <WagmiProvider config={config}>
      <QueryClientProvider client={queryClient}>
        <RainbowKitProvider>
          <ConnectButton.Custom>
            {({
              account,
              chain,
              openAccountModal,
              openChainModal,
              openConnectModal,
              authenticationStatus,
              mounted,
            }) => {
              const ready = mounted && authenticationStatus !== 'loading';
              const connected = ready && account && chain && (!authenticationStatus || authenticationStatus === 'authenticated');
              console.log(account, mounted, authenticationStatus);

              return (
                <div>
                  {(() => {
                    let description: React.ReactNode;
                    let actionList: React.ReactNode;

                    if (!connected) {
                      description = <>You haven't connect wallet yet.</>;
                      actionList = <Button onClick={() => openConnectModal()}>Connect</Button>;
                    } else if (chain.unsupported) {
                      description = <>Current network isn't supported.</>;
                      actionList = <Button onClick={() => openChainModal()}>Switch network</Button>;
                    } else {
                      description = <>You've connected to <strong style={{ cursor: 'pointer' }} onClick={() => openAccountModal()}>{account.address.replace(account.address.slice(6, -4), '...')}</strong>, then you can:</>;
                      actionList = (
                        <>
                          <Button onClick={() => alert('Oops!')}>Authorize</Button>
                          <Button variant="outline-secondary" onClick={() => openAccountModal()}>Change account</Button>
                        </>
                      );
                    }

                    return (
                      <>
                        <p>{description}</p>
                        <div className="d-grid gap-2">{actionList}</div>
                      </>
                    );
                  })()}
                </div>
              );
            }}
          </ConnectButton.Custom>
        </RainbowKitProvider>
      </QueryClientProvider>
    </WagmiProvider>
  );
}

export default WalletAuthorizer;
