package Cx

import data.generic.common as common_lib

CxPolicy[result] {
	resource := input.document[i].resource.aws_iam_user_login_profile[name]

	common_lib.valid_key(resource,"password_reset_required")

    not resource.password_reset_required

	result := {
		"documentId": input.document[i].id,
		"searchKey": sprintf("aws_iam_user_login_profile[%s].password_reset_required", [name]),
		"issueType": "IncorrectValue",
		"keyExpectedValue": "Attribute 'password_reset_required' is true",
		"keyActualValue": "Attribute 'password_reset_required' is false",
	}
}


CxPolicy[result] {
	resource := input.document[i].resource.aws_iam_user_login_profile[name]

	resource.password_length < 14

	result := {
		"documentId": input.document[i].id,
		"searchKey": sprintf("aws_iam_user_login_profile[%s].password_length", [name]),
		"issueType": "IncorrectValue",
		"keyExpectedValue": "Attribute 'password_length' is 14 or grater",
		"keyActualValue": "Attribute 'password_length' is smaller than 14",
	}
}
