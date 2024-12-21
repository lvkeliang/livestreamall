package util

import "errors"

var FieldsError error = errors.New("inquiryError: there was a problem querying the database")

var FormError error = errors.New("The submitted form has an error")

var AreadyLikedError error = errors.New("It is aready liked")

var AreadyCollectedError error = errors.New("It is aready collectd")

var NoArticleExistsError = errors.New("No article exists")

var NoCommectExistsError = errors.New("No commect exists")
