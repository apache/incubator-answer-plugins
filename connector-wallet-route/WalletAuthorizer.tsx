import { ConnectButton } from '@rainbow-me/rainbowkit';
import { useConfig, useSwitchChain, useSignMessage } from 'wagmi';
import { sha256 } from 'js-sha256';
import { useTranslation } from 'react-i18next';
import { Button } from 'react-bootstrap';

function getSearchParamValue(key: string, defaultValue: string = '') {
  return location.search && new URLSearchParams(location.search.slice()).get(key) || defaultValue;
}

function WalletAuthorizer() {
  const { chains } = useConfig();
  const { switchChain } = useSwitchChain();
  const { signMessageAsync } = useSignMessage();

  const { t } = useTranslation('plugin', {
    keyPrefix: 'connector_wallet_route.frontend',
  });

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
          const nonce = getSearchParamValue('nonce', sha256(address));
          const signature = await signMessageAsync({ message: nonce });

          location.href = `/answer/api/v1/connector/redirect/wallet?message=${nonce}&signature=${signature}&address=${address}&redirect=${getSearchParamValue('redirect')}`;
        }

        return (
          <div>
            {(() => {
              let description: React.ReactNode;
              let actionList: React.ReactNode;

              if (!connected) {
                description = <>{t('wallet_needed')}</>;
                actionList = <Button onClick={() => openConnectModal()}>{t('connect_button')}</Button>;
              } else if (chain.unsupported) {
                description = <>{t('wrong_network')}</>;
                actionList = <Button onClick={handleSwitchNetwork}>{t('switch_button')}</Button>;
              } else {
                const translatedArr = t('connected_wallet').split('${ADDRESS}') as React.ReactNode[];
                translatedArr.splice(1, 0, <strong style={{ cursor: 'pointer' }} onClick={() => openAccountModal()}>{address.replace(address.slice(6, -4), '...')}</strong>);
                description = <>{translatedArr.map((c, i) => <span key={i}>{c}</span>)}</>;
                actionList = (
                  <>
                    <Button onClick={handleAuthorize}>{t('authorize_button')}</Button>
                    <Button variant="outline-secondary" onClick={() => openAccountModal()}>{t('disconnect_button')}</Button>
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
