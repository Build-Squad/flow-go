/*
 * Access API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package models

// TransactionExecution : This value indicates whether the transaction execution succeded or not, this value should be checked when determining transaction success.
type TransactionExecution string

// List of TransactionExecution
const (
	PENDING_TransactionExecution TransactionExecution = "Pending"
	SUCCESS_TransactionExecution TransactionExecution = "Success"
	FAILURE_TransactionExecution TransactionExecution = "Failure"
)
