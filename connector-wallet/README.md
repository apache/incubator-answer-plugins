# Wallet connector
> Wallet connector is a OAuth plug-in designed to support Wallet OAuth login.

## How to use

### Build
```bash
./answer build --with github.com/apache/incubator-answer-plugins/connector-wallet
```

### Use Case

- Step 1: Install the wallet plug-in on your chrome/firefox browser. ex:MetaMask,BitgetWallet.

![./imgs/install.png](./imgs/install.png)

- Step 2: Generate your web3 wallet.

![./imgs/create.png](./imgs/create.png)

![./imgs/wallet.png](./imgs/wallet.png)


- Step3: Build your answer with golang and register the plugin.

![./imgs/activate.png](./imgs/activate.png)

- Step4: Log in through your wallet and bind to your email.
![./imgs/click1.png](./imgs/click1.png)

![./imgs/click2.png](./imgs/click2.png)

![./imgs/bind.png](./imgs/bind.png)
