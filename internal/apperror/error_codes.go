package apperror

import "net/http"

var UserErrors = struct {
	InvalidRequest,
	NotFound,
	Conflict,
	Unauthorized,
	FailedToVerify,
	FailedToCreate,
	FailedToLogin,
	FailedToGet,
	FailedToUpdate,
	FailedToLogout,
	Internal *AppError
}{
	InvalidRequest: NewAppError(2001, http.StatusBadRequest, "invalid user request"),
	NotFound:       NewAppError(2002, http.StatusNotFound, "user not found"),
	Conflict:       NewAppError(2003, http.StatusConflict, "user already exists"),
	Unauthorized:   NewAppError(2004, http.StatusUnauthorized, "invalid credentials"),
	FailedToVerify: NewAppError(2101, http.StatusInternalServerError, "failed to verify user"),
	FailedToCreate: NewAppError(2102, http.StatusInternalServerError, "failed to create user"),
	FailedToLogin:  NewAppError(2103, http.StatusInternalServerError, "failed to login user"),
	FailedToGet:    NewAppError(2104, http.StatusInternalServerError, "failed to get user"),
	FailedToUpdate: NewAppError(2105, http.StatusInternalServerError, "failed to update user"),
	FailedToLogout: NewAppError(2106, http.StatusInternalServerError, "failed to logout user"),
	Internal:       NewAppError(2500, http.StatusInternalServerError, "internal server error"),
}

var ContactErrors = struct {
	InvalidRequest,
	NotFound,
	FailedToCreate,
	FailedToUpdate,
	FailedToGet,
	FailedToDelete,
	FailedToSearch,
	Internal *AppError
}{
	InvalidRequest: NewAppError(3001, http.StatusBadRequest, "invalid contact request"),
	NotFound:       NewAppError(3002, http.StatusNotFound, "contact not found"),
	FailedToCreate: NewAppError(3101, http.StatusInternalServerError, "failed to create contact"),
	FailedToUpdate: NewAppError(3102, http.StatusInternalServerError, "failed to update contact"),
	FailedToGet:    NewAppError(3103, http.StatusInternalServerError, "failed to get contact"),
	FailedToDelete: NewAppError(3104, http.StatusInternalServerError, "failed to delete contact"),
	FailedToSearch: NewAppError(3105, http.StatusInternalServerError, "failed to search contacts"),
	Internal:       NewAppError(3500, http.StatusInternalServerError, "internal server error"),
}

var AddressErrors = struct {
	InvalidRequest,
	NotFound,
	FailedToCreate,
	FailedToUpdate,
	FailedToGet,
	FailedToDelete,
	FailedToList,
	Internal *AppError
}{
	InvalidRequest: NewAppError(4001, http.StatusBadRequest, "invalid address request"),
	NotFound:       NewAppError(4002, http.StatusNotFound, "address not found"),
	FailedToCreate: NewAppError(4101, http.StatusInternalServerError, "failed to create address"),
	FailedToUpdate: NewAppError(4102, http.StatusInternalServerError, "failed to update address"),
	FailedToGet:    NewAppError(4103, http.StatusInternalServerError, "failed to get address"),
	FailedToDelete: NewAppError(4104, http.StatusInternalServerError, "failed to delete address"),
	FailedToList:   NewAppError(4105, http.StatusInternalServerError, "failed to list addresses"),
	Internal:       NewAppError(4500, http.StatusInternalServerError, "internal server error"),
}

var AuthErrors = struct {
	Unauthorized,
	MissingToken *AppError
}{
	Unauthorized: NewAppError(5001, http.StatusUnauthorized, "unauthorized"),
	MissingToken: NewAppError(5002, http.StatusUnauthorized, "missing token"),
}

var ItemErrors = struct {
	InvalidRequest,
	FailedToCreateTransaction,
	FailedToCreateItem *AppError
}{
	InvalidRequest:            NewAppError(6001, http.StatusBadRequest, "invalid item request"),
	FailedToCreateTransaction: NewAppError(6002, http.StatusBadRequest, "faild to create transaction to create item"),
	FailedToCreateItem:        NewAppError(6003, http.StatusBadRequest, "faild to create item"),
}

var UnknownErrors = struct {
	SQLNoRowsError,
	UnknownValidationError *AppError
}{
	SQLNoRowsError:         NewAppError(99091, http.StatusNotFound, "qrm: no rows in result set"),
	UnknownValidationError: NewAppError(99090, http.StatusBadRequest, "unknown validation error"),
}
