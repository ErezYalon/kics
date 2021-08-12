package Cx

import data.generic.common as common_lib

CxPolicy[result] {
	resource := input.document[i].resource.aws_db_instance[name]
	common_lib.valid_key(resource,"auto_minor_version_upgrade")
    not resource.auto_minor_version_upgrade

	result := {
		"documentId": input.document[i].id,
		"searchKey": sprintf("aws_db_instance[%s].auto_minor_version_upgrade", [name]),
		"issueType": "IncorrectValue",
		"keyExpectedValue": "'aws_db_instance.auto_minor_version_upgrade'  is 'true'",
		"keyActualValue": "'aws_db_instance.auto_minor_version_upgrade'  is 'false'",
	}
}

