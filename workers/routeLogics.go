package workers

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"miniwallet.co.id/models"
)

/**
DoInitWallet is used to create an account as well as getting the token for the other API endpoints.
**/
func DoInitWallet(w http.ResponseWriter, r *http.Request) {
	payloadValidator := ValidateRequestBody(r)

	if !payloadValidator {
		apiResponse := PrepareAPIResponse(false, "Missing data for required field.", false)
		sendAPIResponse(apiResponse, w)
	} else {

		customers := readUserJson(w)
		// Get customer_xid from Form Data
		customer_id := r.FormValue("customer_xid")

		// Check if the customer_id exist in customer data
		counter := 0
		for _, v := range customers {
			if v.CustomerID == customer_id {
				var varToken models.TokenResponse
				varToken.Token = v.TokenID
				apiResponse := PrepareAPIResponse(true, "", varToken)
				sendAPIResponse(apiResponse, w)
				return
			}
		}

		// Conditions if customer id not found in Database
		if counter < 1 {
			apiResponse := PrepareAPIResponse(false, "Customer ID Not Found!", false)
			sendAPIResponse(apiResponse, w)
		}
	}
}

/**
DoEnableWallet is used to enable the wallet to stores the virtual money that customers can apply for approval.
**/
func DoEnableWallet(w http.ResponseWriter, r *http.Request) {

	// Get the user data from User JSON
	customers := readUserJson(w)
	var customerID string

	/**
		Token Validation
		Get the Customer ID from TokenId given in request body
	**/
	tokenValidation := false
	for _, v := range customers {
		validateToken := ValidateToken(r, v.TokenID)
		if validateToken {
			customerID = v.CustomerID
			tokenValidation = true
			break
		}
	}

	if !tokenValidation {
		apiResponse := PrepareAPIResponse(false, "UnAuthorized", false)
		sendAPIResponse(apiResponse, w)
		return
	}

	/**
		Read Wallet data from wallet Json file
		Get the wallet data based on customer id
	**/
	wallet := readWallet(w)
	for i := 0; i < len(wallet); i++ {
		if wallet[i].OwnedBy == customerID {
			if wallet[i].Status == "enabled" {
				apiResponse := PrepareAPIResponse(false, "Already Enabled", false)
				sendAPIResponse(apiResponse, w)
				return
			}
			if wallet[i].Status == "disabled" {
				wallet[i].Status = "enabled"
				wallet[i].EnabledAt = time.Now().String()
				updateWalletData(wallet, w)
				apiResponse := PrepareAPIResponse(true, "", wallet[i])
				sendAPIResponse(apiResponse, w)
			}
		}
	}
}

/**
DoDisableWallet is used to disabled the wallet
**/
func DoDisableWallet(w http.ResponseWriter, r *http.Request) {

	// Get the user data from User JSON
	customers := readUserJson(w)
	var customerID string

	/**
		Token Validation
		Get the Customer ID from TokenId given in request body
	**/
	tokenValidation := false
	for _, v := range customers {
		validateToken := ValidateToken(r, v.TokenID)
		if validateToken {
			customerID = v.CustomerID
			tokenValidation = true
			break
		}
	}

	if !tokenValidation {
		apiResponse := PrepareAPIResponse(false, "UnAuthorized", false)
		sendAPIResponse(apiResponse, w)
		return
	}

	/**
		Read Wallet data from wallet Json file
		Get the wallet data based on customer id
	**/
	wallet := readWallet(w)
	for i := 0; i < len(wallet); i++ {
		if wallet[i].OwnedBy == customerID {
			if wallet[i].Status == "disabled" {
				apiResponse := PrepareAPIResponse(false, "Already Disabled", false)
				sendAPIResponse(apiResponse, w)
				return
			}
			if wallet[i].Status == "enabled" {
				wallet[i].Status = "disabled"
				wallet[i].DisabledAt = time.Now().String()
				updateWalletData(wallet, w)
				apiResponse := PrepareAPIResponse(true, "", wallet[i])
				sendAPIResponse(apiResponse, w)
			}
		}
	}
}

/**
ViewWallet is used to view the wallet information
**/
func ViewWallet(w http.ResponseWriter, r *http.Request) {
	// Get the user data from User JSON
	customers := readUserJson(w)
	var customerID string

	/**
		Token Validation
		Get the Customer ID from TokenId given in request body
	**/
	tokenValidation := false
	for _, v := range customers {
		validateToken := ValidateToken(r, v.TokenID)
		if validateToken {
			customerID = v.CustomerID
			tokenValidation = true
			break
		}
	}

	if !tokenValidation {
		apiResponse := PrepareAPIResponse(false, "UnAuthorized", false)
		sendAPIResponse(apiResponse, w)
		return
	}

	/**
		Read Wallet data from wallet Json file
		Get the wallet data based on customer id
	**/
	wallet := readWallet(w)
	for i := 0; i < len(wallet); i++ {
		if wallet[i].OwnedBy == customerID {
			if wallet[i].Status == "disabled" {
				apiResponse := PrepareAPIResponse(false, "Disabled", false)
				sendAPIResponse(apiResponse, w)
				return
			}
			apiResponse := PrepareAPIResponse(true, "", wallet[i])
			sendAPIResponse(apiResponse, w)
			return
		}
	}
	apiResponse := PrepareAPIResponse(false, "You don't have Wallet", false)
	sendAPIResponse(apiResponse, w)
}

/**
DepositMoney is used to add virtual money into wallet
**/
func DepositMoney(w http.ResponseWriter, r *http.Request) {

	/**
		Token Validation
		Get the Customer ID from TokenId given in request body
	**/
	customers := readUserJson(w)
	var customerID string

	tokenValidation := false
	for _, v := range customers {
		validateToken := ValidateToken(r, v.TokenID)
		if validateToken {
			customerID = v.CustomerID
			tokenValidation = true
			break
		}
	}

	if !tokenValidation {
		apiResponse := PrepareAPIResponse(false, "UnAuthorized", false)
		sendAPIResponse(apiResponse, w)
		return
	}

	/**
		Request Body Validation
		Amount and Reference id mandatory in Request param
	**/
	payloadValidator := ValidateRequestBody(r)
	if !payloadValidator {
		apiResponse := PrepareAPIResponse(false, "Missing data for required field.", false)
		sendAPIResponse(apiResponse, w)
	}

	amount := r.FormValue("amount")
	referenceId := r.FormValue("reference_id")

	if amount == "" {
		apiResponse := PrepareAPIResponse(false, "Missing amount for required field.", false)
		sendAPIResponse(apiResponse, w)
	}

	if referenceId == "" {
		apiResponse := PrepareAPIResponse(false, "Missing reference_id for required field.", false)
		sendAPIResponse(apiResponse, w)
	}

	/**
		Processing the deposit:
		1. Check if the Reference ID exist or not
		2. Preparing new deposit data from Payload
		3. Add new deposit data into database (deposit.json)
		4. Update amount in wallet
	**/
	var deposits []models.Deposit
	var amountInt int
	amountInt, _ = strconv.Atoi(amount)
	deposits = readDeposit(w)
	if deposits == nil {
		return
	}

	// Validate the reference id exist in database or not
	var i = 0
	if len(deposits) > 0 {
		for i = 0; i < len(deposits); i++ {
			if deposits[i].ReferenceID == referenceId {
				apiResponse := PrepareAPIResponse(false, "Reference ID is exist", false)
				sendAPIResponse(apiResponse, w)
				return
			}
		}
	}

	// Prepare new Deposit Data
	var deposit models.Deposit
	var depositID string
	depositID = "deposit" + time.Now().String()
	depositID = base64.StdEncoding.EncodeToString([]byte(depositID))
	deposit.ID = depositID
	deposit.DepositBy = customerID
	deposit.DepositAt = time.Now().String()
	deposit.Amount = amountInt
	deposit.ReferenceID = referenceId
	deposit.Status = "success"
	deposits = append(deposits, deposit)

	// Read and Update the wallet data
	wallet := readWallet(w)

	for i := 0; i < len(wallet); i++ {
		if wallet[i].OwnedBy == customerID {
			if wallet[i].Status == "enabled" {
				amountInt, _ = strconv.Atoi(amount)
				wallet[i].Balance = amountInt

				// Save new deposit data
				updateDepositData(deposits, w)
				// Update amount in the wallet
				updateWalletData(wallet, w)
				break
			}
			if wallet[i].Status == "disabled" {
				apiResponse := PrepareAPIResponse(false, "Disabled", false)
				sendAPIResponse(apiResponse, w)
				return
			}
		}
	}

	apiResponse := PrepareAPIResponse(true, "", deposit)
	sendAPIResponse(apiResponse, w)
	return
}

/**
WithdrawMoney is used to reduce / withdraw virtual money into wallet
**/
func WithdrawMoney(w http.ResponseWriter, r *http.Request) {

	/**
		Token Validation
		Get the Customer ID from TokenId given in request body
	**/
	customers := readUserJson(w)
	var customerID string

	tokenValidation := false
	for _, v := range customers {
		validateToken := ValidateToken(r, v.TokenID)
		if validateToken {
			customerID = v.CustomerID
			tokenValidation = true
			break
		}
	}

	if !tokenValidation {
		apiResponse := PrepareAPIResponse(false, "UnAuthorized", false)
		sendAPIResponse(apiResponse, w)
		return
	}

	/**
		Request Body Validation
		Amount and Reference id mandatory in Request param
	**/
	payloadValidator := ValidateRequestBody(r)
	if !payloadValidator {
		apiResponse := PrepareAPIResponse(false, "Missing data for required field.", false)
		sendAPIResponse(apiResponse, w)
	}

	amount := r.FormValue("amount")
	referenceId := r.FormValue("reference_id")

	if amount == "" {
		apiResponse := PrepareAPIResponse(false, "Missing amount for required field.", false)
		sendAPIResponse(apiResponse, w)
	}

	if referenceId == "" {
		apiResponse := PrepareAPIResponse(false, "Missing reference_id for required field.", false)
		sendAPIResponse(apiResponse, w)
	}

	var withdrawals []models.Withdraw
	var amountInt int
	amountInt, _ = strconv.Atoi(amount)
	withdrawals = readDWithdraw(w)

	if withdrawals == nil {
		return
	}
	var i = 0

	// Validate the reference id exist in database or not
	if len(withdrawals) > 0 {
		for i = 0; i < len(withdrawals); i++ {
			if withdrawals[i].ReferenceID == referenceId {
				apiResponse := PrepareAPIResponse(false, "Reference ID is exist", false)
				sendAPIResponse(apiResponse, w)
				return
			}
		}
	}

	/**
		Processing the withdrawal:
		1. Check if the Reference ID exist or not
		2. Preparing new withdrawal data from Payload
		3. Add new withdrawal data into database (withdrawal.json)
		4. Update amount in wallet
	**/

	// Prepare new Withdrawal Data
	var withdraw models.Withdraw
	var withdrawID string
	withdrawID = "withdraw" + time.Now().String()
	withdrawID = base64.StdEncoding.EncodeToString([]byte(withdrawID))

	withdraw.ID = withdrawID
	withdraw.WithdrawnBy = customerID
	withdraw.WithdrawnAt = time.Now().String()
	withdraw.Amount = amountInt
	withdraw.ReferenceID = referenceId
	withdraw.Status = "success"
	withdrawals = append(withdrawals, withdraw)

	// Update the Amount in the wallet
	// Validation of amount in wallet before do withdrawal
	wallet := readWallet(w)

	for i := 0; i < len(wallet); i++ {
		if wallet[i].OwnedBy == customerID {
			if wallet[i].Status == "enabled" {
				amountInt, _ = strconv.Atoi(amount)
				if amountInt > wallet[i].Balance {
					apiResponse := PrepareAPIResponse(false, "Unsufficient Amount in your wallet", false)
					sendAPIResponse(apiResponse, w)
					return
				}
				wallet[i].Balance = wallet[i].Balance - amountInt

				// Save new withdrawal data
				updateWithdrawData(withdrawals, w)
				// Update amount in the wallet
				updateWalletData(wallet, w)
				break
			}
			if wallet[i].Status == "disabled" {
				apiResponse := PrepareAPIResponse(false, "Disabled", false)
				sendAPIResponse(apiResponse, w)
				return
			}
		}
	}

	apiResponse := PrepareAPIResponse(true, "", withdraw)
	sendAPIResponse(apiResponse, w)
	return
}

/**
readUserJson is used to read the data from user.json that act as database
**/
func readUserJson(w http.ResponseWriter) models.CustomerData {
	// Validate and Read the User Json File
	validationStatus, jsonData := ValidateJSONFileReader("./database/user.json")

	// If Error when reading JSON file
	if !validationStatus {
		apiResponse := PrepareAPIResponse(false, "Error when Reading JSON File", false)
		sendAPIResponse(apiResponse, w)
	}

	// If Succeed to read JSON File
	var customers models.CustomerData
	if validationStatus {
		err := json.Unmarshal(jsonData, &customers)
		if err != nil {
			apiResponse := PrepareAPIResponse(false, "Error when assign value", false)
			sendAPIResponse(apiResponse, w)
		}
	}

	return customers
}

/**
readWallet is used to read the data from wallet.json that act as database
**/
func readWallet(w http.ResponseWriter) models.Wallet {
	// Validate and Read the User Json File
	validationStatus, jsonData := ValidateJSONFileReader("./database/wallet.json")

	// If Error when reading JSON file
	if !validationStatus {
		apiResponse := PrepareAPIResponse(false, "Error when Reading JSON File", false)
		sendAPIResponse(apiResponse, w)
	}

	// If Succeed to read JSON File
	var wallet models.Wallet
	if validationStatus {
		err := json.Unmarshal(jsonData, &wallet)
		if err != nil {
			apiResponse := PrepareAPIResponse(false, "Error when assign value", false)
			sendAPIResponse(apiResponse, w)
		}
	}

	return wallet
}

/**
readDeposit is used to read the data from deposit.json that act as database
**/
func readDeposit(w http.ResponseWriter) []models.Deposit {
	// Validate and Read the User Json File
	validationStatus, jsonData := ValidateJSONFileReader("./database/deposit.json")

	// If Error when reading JSON file
	if !validationStatus {
		apiResponse := PrepareAPIResponse(false, "Error when Reading JSON File", false)
		sendAPIResponse(apiResponse, w)
	}

	// If Succeed to read JSON File
	var deposits []models.Deposit

	if validationStatus && len(jsonData) > 0 {
		err := json.Unmarshal(jsonData, &deposits)
		if err != nil {
			apiResponse := PrepareAPIResponse(false, "Error when assign value", false)
			sendAPIResponse(apiResponse, w)
		}
	}

	return deposits
}

/**
readDWithdraw is used to read the data from withdrawal.json that act as database
**/
func readDWithdraw(w http.ResponseWriter) []models.Withdraw {
	// Validate and Read the User Json File
	validationStatus, jsonData := ValidateJSONFileReader("./database/withdrawal.json")

	// If Error when reading JSON file
	if !validationStatus {
		apiResponse := PrepareAPIResponse(false, "Error when Reading JSON File", false)
		sendAPIResponse(apiResponse, w)
		return nil
	}

	// If Succeed to read JSON File
	var withdrawals []models.Withdraw

	if validationStatus && len(jsonData) > 0 {
		err := json.Unmarshal(jsonData, &withdrawals)
		if err != nil {
			apiResponse := PrepareAPIResponse(false, "Error when assign value", false)
			sendAPIResponse(apiResponse, w)
		}
	}

	return withdrawals
}

/**
updateWalletData is used to update the data in wallet.json after the process of deposit or withdrawal
**/
func updateWalletData(wallet models.Wallet, w http.ResponseWriter) error {
	walletDataJson, err := json.Marshal(wallet)

	if err != nil {
		apiResponse := PrepareAPIResponse(false, "Error when updating wallet", false)
		sendAPIResponse(apiResponse, w)
	}

	err = ioutil.WriteFile("./database/wallet.json", walletDataJson, 0644)
	return err
}

/**
updateDepositData is used to update or add the data in deposit.json after the process of deposit
**/
func updateDepositData(deposits []models.Deposit, w http.ResponseWriter) error {

	depositsDataJson, err := json.Marshal(deposits)
	if err != nil {
		apiResponse := PrepareAPIResponse(false, "Error when updating wallet", false)
		sendAPIResponse(apiResponse, w)
	}

	err = ioutil.WriteFile("./database/deposit.json", depositsDataJson, 0644)
	return err
}

/**
updateWithdrawData is used to update or add the data in withdrawal.json after the process of withdrawal
**/
func updateWithdrawData(withdrawals []models.Withdraw, w http.ResponseWriter) error {

	withdrawalsDataJson, err := json.Marshal(withdrawals)
	if err != nil {
		apiResponse := PrepareAPIResponse(false, "Error when updating wallet", false)
		sendAPIResponse(apiResponse, w)
	}

	err = ioutil.WriteFile("./database/withdrawals.json", withdrawalsDataJson, 0644)
	return err
}
