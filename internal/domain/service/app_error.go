package service // ошибки сервисного слоя

type Code string // тип для кодов ошибок

const (
	CodeValidation Code = "validation_error" // неверные данные
	CodeNotFound   Code = "not_found"        // не найдено
	CodeInternal   Code = "internal_error"   // внутренняя ошибка
)

type AppError struct { // единый тип ошибки сервиса
	Code    Code  // код ошибки
	Details any   // детали (опц.)
	err     error // исходная ошибка (опц.)
}

func (e *AppError) Error() string { // текст ошибки
	if e.err != nil { // если есть причина
		return string(e.Code) + ": " + e.err.Error() // код + причина
	}
	return string(e.Code) // только код
}

func (e *AppError) Unwrap() error { return e.err } // для errors.Is/As

// helpers
func Validation(details any) error { return &AppError{Code: CodeValidation, Details: details} } // создать validation
func NotFound(details any) error   { return &AppError{Code: CodeNotFound, Details: details} }   // создать not_found
func Internal(err error) error     { return &AppError{Code: CodeInternal, err: err} }           // создать internal
