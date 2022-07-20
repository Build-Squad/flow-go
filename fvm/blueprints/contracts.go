package blueprints

import (
	_ "embed"
	"encoding/hex"
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/flow-go/fvm/utils"
	"github.com/onflow/flow-go/model/flow"
)

var ContractDeploymentAuthorizedAddressesPath = cadence.Path{
	Domain:     common.PathDomainStorage.Identifier(),
	Identifier: "authorizedAddressesToDeployContracts",
}
var ContractRemovalAuthorizedAddressesPath = cadence.Path{
	Domain:     common.PathDomainStorage.Identifier(),
	Identifier: "authorizedAddressesToRemoveContracts",
}
var IsContractDeploymentRestrictedPath = cadence.Path{
	Domain:     common.PathDomainStorage.Identifier(),
	Identifier: "isContractDeploymentRestricted",
}

//go:embed scripts/setContractOperationAuthorizersTransactionTemplate.cdc
var setContractOperationAuthorizersTransactionTemplate string

//go:embed scripts/setIsContractDeploymentRestrictedTransactionTemplate.cdc
var setIsContractDeploymentRestrictedTransactionTemplate string

//go:embed scripts/deployContractTransactionTemplate.cdc
var deployContractTransactionTemplate string

// SetContractDeploymentAuthorizersTransaction returns a transaction for updating list of authorized accounts allowed to deploy/update contracts
func SetContractDeploymentAuthorizersTransaction(serviceAccount flow.Address, authorized []flow.Address) (*flow.TransactionBody, error) {
	return setContractAuthorizersTransaction(ContractDeploymentAuthorizedAddressesPath, serviceAccount, authorized)
}

// SetContractRemovalAuthorizersTransaction returns a transaction for updating list of authorized accounts allowed to remove contracts
func SetContractRemovalAuthorizersTransaction(serviceAccount flow.Address, authorized []flow.Address) (*flow.TransactionBody, error) {
	return setContractAuthorizersTransaction(ContractRemovalAuthorizedAddressesPath, serviceAccount, authorized)
}

func setContractAuthorizersTransaction(
	path cadence.Path,
	serviceAccount flow.Address,
	authorized []flow.Address,
) (*flow.TransactionBody, error) {
	addresses := utils.FlowAddressSliceToCadenceAddressSlice(authorized)
	addressesArg, err := jsoncdc.Encode(utils.AddressSliceToCadenceValue(addresses))
	if err != nil {
		return nil, err
	}

	pathArg, err := jsoncdc.Encode(path)
	if err != nil {
		return nil, err
	}

	return flow.NewTransactionBody().
		SetScript([]byte(setContractOperationAuthorizersTransactionTemplate)).
		AddAuthorizer(serviceAccount).
		AddArgument(addressesArg).
		AddArgument(pathArg), nil
}

// SetIsContractDeploymentRestrictedTransaction sets the restricted flag for contract deployment
func SetIsContractDeploymentRestrictedTransaction(serviceAccount flow.Address, restricted bool) (*flow.TransactionBody, error) {
	argRestricted, err := jsoncdc.Encode(cadence.Bool(restricted))
	if err != nil {
		return nil, err
	}

	argPath, err := jsoncdc.Encode(IsContractDeploymentRestrictedPath)
	if err != nil {
		return nil, err
	}

	return flow.NewTransactionBody().
		SetScript([]byte(setIsContractDeploymentRestrictedTransactionTemplate)).
		AddAuthorizer(serviceAccount).
		AddArgument(argRestricted).
		AddArgument(argPath), nil
}

// TODO (ramtin) get rid of authorizers
func DeployContractTransaction(address flow.Address, contract []byte, contractName string) *flow.TransactionBody {
	return flow.NewTransactionBody().
		SetScript([]byte(deployContractTransactionTemplate)).
		AddArgument(jsoncdc.MustEncode(cadence.String(contractName))).
		AddArgument(jsoncdc.MustEncode(cadence.String(hex.EncodeToString(contract)))).
		AddAuthorizer(address)
}
