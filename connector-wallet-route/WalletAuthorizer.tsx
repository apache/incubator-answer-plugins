import '@rainbow-me/rainbowkit/styles.css';
import { ConnectButton } from '@rainbow-me/rainbowkit';
import { useConfig, useSwitchChain, useSignMessage } from 'wagmi';
import { sha256 } from 'js-sha256';
// import { useTranslation } from 'react-i18next';
import { Button } from 'react-bootstrap';

function resolveNonce(address: string) {
  return location.search && new URLSearchParams(location.search.slice()).get('nonce') || sha256(address);
}

function WalletAuthorizer() {
  const { chains } = useConfig();
  const { switchChain } = useSwitchChain();
  const { signMessageAsync } = useSignMessage();

  // const { t } = useTranslation('plugin', {
  //   keyPrefix: 'connector_wallet_route.frontend',
  // });

  return (
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
        const { address } = account || { address: '' };

        const handleSwitchNetwork = () => {
          if (chains.length === 1) {
            switchChain({ chainId: chains[0].id });
          } else {
            openChainModal();
          }
        };

        const handleAuthorize = async () => {
          const nonce = resolveNonce(address);
          const signature = await signMessageAsync({ message: nonce });

          location.href = `/answer/api/v1/connector/redirect/wallet?message=${nonce}&&signature=${signature}&&address=${address}`;
        }

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
                actionList = <Button onClick={handleSwitchNetwork}>Switch network</Button>;
              } else {
                description = <>You've connected to <strong style={{ cursor: 'pointer' }} onClick={() => openAccountModal()}>{address.replace(address.slice(6, -4), '...')}</strong>, then you can:</>;
                actionList = (
                  <>
                    <Button onClick={handleAuthorize}>Authorize</Button>
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
  );
}

export default WalletAuthorizer;
