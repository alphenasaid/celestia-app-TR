# celestia-app

**celestia-app** Cosmos SDK ve [celestia-core] kullanılarak oluşturulmuş bir blok zinciri uygulamasıdır(https://github.com/celestiaorg/celestia-core) in place of tendermint.sorumluluk reddi: **WIP**

## kurmak
```
make install
```

### Kendi tek düğümlü devnet'inizi oluşturun
```
celestia-appd init test --chain-id test
celestia-appd keys add user1
celestia-appd add-genesis-account <address from above command> 10000000utia,1000token
celestia-appd gentx user1 1000000utia --chain-id test
celestia-appd collect-gentxs
celestia-appd start
```
## kullanma
Yerel bir olaya veri göndermek için 'celestia-appd' arka plan programı cli komutunu kullanın
  
```celestia-appd tx payment payForData [hexNamespace] [hexMessage] [flags]```
