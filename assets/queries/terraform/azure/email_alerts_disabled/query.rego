package Cx

CxPolicy[result] {
	resource := input.document[i].resource.azurerm_security_center_contact[name]

	object.get(resource, "alert_notifications", "undefined") != "undefined"
	not resource.alert_notifications

	result := {
		"documentId": input.document[i].id,
		"searchKey": sprintf("azurerm_security_center_contact[%s].alert_notifications", [name]),
		"issueType": "IncorrectValue",
		"keyExpectedValue": sprintf("'azurerm_security_center_contact.%s.alert_notifications' is true", [name]),
		"keyActualValue": sprintf("'azurerm_security_center_contact.%s.alert_notifications' is false", [name]),
	}
}
