package correspondence

// This file contains predefined templates for common police correspondence

// GetDefaultTemplates returns a list of default correspondence templates
func GetDefaultTemplates() []*Template {
	return []*Template{
		// Evidence request template
		{
			ID:      "TMPL-EVIDENCE-REQ-1",
			Name:    "Standard Evidence Request",
			Type:    TypeEvidenceReq,
			Subject: "Evidence Request: Case {{.CaseNumber}}",
			Body: `
To: {{.RecipientName}}
{{.RecipientOrganization}}
{{.RecipientAddress}}

RE: Request for Evidence - Case #{{.CaseNumber}}

Dear {{.RecipientTitle}} {{.RecipientLastName}},

This letter serves as a formal request for evidence related to Case #{{.CaseNumber}}, involving {{.CaseDescription}}.

The {{.DepartmentName}} is investigating this matter and requires access to the following items:

{{.EvidenceList}}

These items are being requested pursuant to {{.LegalAuthority}} and are essential to our ongoing investigation. Please provide these items by {{.Deadline}}.

If you have any questions or concerns regarding this request, please contact me directly at {{.ContactPhone}} or via email at {{.ContactEmail}}.

Your prompt attention to this matter is greatly appreciated.

Sincerely,

{{.OfficerName}}
{{.OfficerTitle}}, Badge #{{.BadgeNumber}}
{{.DepartmentName}}
`,
			TemplateVars: []string{
				"CaseNumber", "RecipientName", "RecipientOrganization", "RecipientAddress",
				"RecipientTitle", "RecipientLastName", "CaseDescription", "DepartmentName",
				"EvidenceList", "LegalAuthority", "Deadline", "ContactPhone", "ContactEmail",
				"OfficerName", "OfficerTitle", "BadgeNumber",
			},
			Department: "Any",
			IsApproved: true,
		},

		// Witness interview request
		{
			ID:      "TMPL-INTERVIEW-REQ-1",
			Name:    "Witness Interview Request",
			Type:    TypeLetter,
			Subject: "Request for Interview: Case {{.CaseNumber}}",
			Body: `
To: {{.RecipientName}}
{{.RecipientAddress}}

RE: Request for Interview - Case #{{.CaseNumber}}

Dear {{.RecipientTitle}} {{.RecipientLastName}},

I am writing to request your cooperation in an ongoing investigation being conducted by the {{.DepartmentName}}. Based on information gathered, we believe you may have witnessed events related to Case #{{.CaseNumber}}, which involves {{.CaseDescription}}.

We would like to schedule an interview with you at your earliest convenience to discuss any information you may have regarding this matter. Your assistance is greatly valued and will help us in conducting a thorough investigation.

The interview can be conducted at {{.InterviewLocation}} or at another location of your choosing. Please contact me at {{.ContactPhone}} or {{.ContactEmail}} to arrange a suitable time.

Please note that this is not an indication that you are under suspicion or involved in any wrongdoing. Your role as a potential witness is important to our fact-finding process.

Thank you for your cooperation in this matter.

Sincerely,

{{.OfficerName}}
{{.OfficerTitle}}, Badge #{{.BadgeNumber}}
{{.DepartmentName}}
`,
			TemplateVars: []string{
				"CaseNumber", "RecipientName", "RecipientAddress", "RecipientTitle",
				"RecipientLastName", "CaseDescription", "DepartmentName", "InterviewLocation",
				"ContactPhone", "ContactEmail", "OfficerName", "OfficerTitle", "BadgeNumber",
			},
			Department: "Any",
			IsApproved: true,
		},

		// Subpoena template
		{
			ID:      "TMPL-SUBPOENA-1",
			Name:    "Standard Subpoena",
			Type:    TypeSubpoena,
			Subject: "Subpoena: Case {{.CaseNumber}}",
			Body: `
STATE OF {{.State}}
COUNTY OF {{.County}}
{{.CourtName}}

SUBPOENA {{.SubpoenaType}}

CASE NUMBER: {{.CaseNumber}}
CASE NAME: {{.CaseName}}

TO: {{.RecipientName}}
    {{.RecipientAddress}}

YOU ARE HEREBY COMMANDED to appear in the {{.CourtName}} at {{.CourtAddress}}, on {{.AppearanceDate}} at {{.AppearanceTime}}, to testify in the above case.

{{if eq .SubpoenaType "DUCES TECUM"}}
YOU ARE ALSO COMMANDED to bring with you the following items:

{{.ItemsList}}
{{end}}

FAILURE TO APPEAR IN ACCORDANCE WITH THIS SUBPOENA MAY BE DEEMED A CONTEMPT OF COURT FOR WHICH YOU MAY BE PUNISHED AS PROVIDED BY LAW.

ISSUED ON: {{.IssueDate}}

BY ORDER OF THE COURT:

{{.JudgeName}}
{{.JudgeTitle}}

REQUESTING OFFICER:
{{.OfficerName}}, {{.OfficerTitle}}
{{.DepartmentName}}
Badge #{{.BadgeNumber}}
Contact: {{.ContactPhone}}
`,
			TemplateVars: []string{
				"State", "County", "CourtName", "SubpoenaType", "CaseNumber", "CaseName",
				"RecipientName", "RecipientAddress", "CourtAddress", "AppearanceDate",
				"AppearanceTime", "ItemsList", "IssueDate", "JudgeName", "JudgeTitle",
				"OfficerName", "OfficerTitle", "DepartmentName", "BadgeNumber", "ContactPhone",
			},
			Department: "Any",
			IsApproved: true,
		},

		// Search warrant template
		{
			ID:      "TMPL-WARRANT-1",
			Name:    "Search Warrant Application",
			Type:    TypeWarrant,
			Subject: "Application for Search Warrant: Case {{.CaseNumber}}",
			Body: `
STATE OF {{.State}}
COUNTY OF {{.County}}
{{.CourtName}}

APPLICATION FOR SEARCH WARRANT

CASE NUMBER: {{.CaseNumber}}

AFFIDAVIT

I, {{.OfficerName}}, Badge #{{.BadgeNumber}}, being duly sworn, depose and say that I have reason to believe that on the premises known as:

{{.PremisesAddress}}
{{.PremisesDescription}}

in the City of {{.City}}, County of {{.County}}, State of {{.State}}, there is now being concealed certain property, namely:

{{.PropertyDescription}}

which is {{.LegalBasis}}

The facts to support a finding of Probable Cause are as follows:

{{.ProbableCauseStatement}}

Wherefore, I request that a Search Warrant be issued authorizing a search of the above-described premises and the seizure of the above-described items.

{{.OfficerName}}, {{.OfficerTitle}}
{{.DepartmentName}}
Badge #{{.BadgeNumber}}

Sworn to before me and subscribed in my presence on {{.SwornDate}}

{{.JudgeName}}
{{.JudgeTitle}}
`,
			TemplateVars: []string{
				"State", "County", "CourtName", "CaseNumber", "OfficerName", "BadgeNumber",
				"PremisesAddress", "PremisesDescription", "City", "PropertyDescription",
				"LegalBasis", "ProbableCauseStatement", "OfficerTitle", "DepartmentName",
				"SwornDate", "JudgeName", "JudgeTitle",
			},
			Department: "Any",
			IsApproved: true,
		},

		// Internal case memo
		{
			ID:      "TMPL-MEMO-1",
			Name:    "Internal Case Memo",
			Type:    TypeMemo,
			Subject: "Case Update: {{.CaseNumber}} - {{.CaseTitle}}",
			Body: `
MEMORANDUM

TO: {{.RecipientName}}, {{.RecipientTitle}}
FROM: {{.OfficerName}}, {{.OfficerTitle}}, Badge #{{.BadgeNumber}}
DATE: {{.CurrentDate}}
RE: Case Update - #{{.CaseNumber}} ({{.CaseTitle}})

CLASSIFICATION: {{.Classification}}

SUMMARY:
This memo provides an update on the investigation of Case #{{.CaseNumber}}, involving {{.CaseDescription}}.

RECENT DEVELOPMENTS:
{{.RecentDevelopments}}

CURRENT STATUS:
{{.CurrentStatus}}

NEXT STEPS:
{{.NextSteps}}

RESOURCES NEEDED:
{{.ResourcesNeeded}}

TIMELINE:
{{.Timeline}}

Please contact me at {{.ContactInfo}} if you require additional information or wish to discuss this case further.

{{.OfficerName}}
{{.OfficerTitle}}
Badge #{{.BadgeNumber}}
{{.DepartmentName}}
`,
			TemplateVars: []string{
				"CaseNumber", "CaseTitle", "RecipientName", "RecipientTitle", "OfficerName",
				"OfficerTitle", "BadgeNumber", "CurrentDate", "CaseDescription", "Classification",
				"RecentDevelopments", "CurrentStatus", "NextSteps", "ResourcesNeeded",
				"Timeline", "ContactInfo", "DepartmentName",
			},
			Department: "Any",
			IsApproved: true,
		},

		// Press release
		{
			ID:      "TMPL-PRESS-1",
			Name:    "Standard Press Release",
			Type:    TypePressRelease,
			Subject: "{{.DepartmentName}} Press Release: {{.Title}}",
			Body: `
PRESS RELEASE
{{.DepartmentName}}
{{.DepartmentAddress}}
{{.DepartmentPhone}}
{{.DepartmentWebsite}}

FOR IMMEDIATE RELEASE
{{.ReleaseDate}}

{{.Title}}

{{.City}}, {{.State}} - {{.Summary}}

{{.BodyParagraph1}}

{{.BodyParagraph2}}

{{.BodyParagraph3}}

{{if .QuotePerson}}
"{{.Quote}}" said {{.QuotePersonTitle}} {{.QuotePerson}}.
{{end}}

{{if .InvestigationStatus}}
INVESTIGATION STATUS:
{{.InvestigationStatus}}
{{end}}

{{if .PublicAssistance}}
REQUEST FOR PUBLIC ASSISTANCE:
{{.PublicAssistance}}
{{end}}

For more information, please contact:
{{.ContactName}}
{{.ContactTitle}}
{{.ContactPhone}}
{{.ContactEmail}}

###
`,
			TemplateVars: []string{
				"DepartmentName", "DepartmentAddress", "DepartmentPhone", "DepartmentWebsite",
				"ReleaseDate", "Title", "City", "State", "Summary", "BodyParagraph1",
				"BodyParagraph2", "BodyParagraph3", "Quote", "QuotePerson", "QuotePersonTitle",
				"InvestigationStatus", "PublicAssistance", "ContactName", "ContactTitle",
				"ContactPhone", "ContactEmail",
			},
			Department: "Public Relations",
			IsApproved: true,
		},

		// Evidence chain of custody form
		{
			ID:      "TMPL-EVIDENCE-CUSTODY-1",
			Name:    "Evidence Chain of Custody Form",
			Type:    TypeEvidence,
			Subject: "Chain of Custody: Evidence #{{.EvidenceNumber}} - Case {{.CaseNumber}}",
			Body: `
{{.DepartmentName}}
EVIDENCE CHAIN OF CUSTODY FORM

CASE NUMBER: {{.CaseNumber}}
EVIDENCE NUMBER: {{.EvidenceNumber}}
EVIDENCE DESCRIPTION: {{.EvidenceDescription}}

RECOVERED BY: {{.RecoveredBy}}, Badge #{{.RecoveredByBadge}}
RECOVERY LOCATION: {{.RecoveryLocation}}
RECOVERY DATE/TIME: {{.RecoveryDateTime}}
RECOVERY NOTES: {{.RecoveryNotes}}

CHAIN OF CUSTODY:

1. FROM: {{.RecoveredBy}}, Badge #{{.RecoveredByBadge}}
   TO: {{.FirstCustodian}}, {{.FirstCustodianTitle}}
   DATE/TIME: {{.FirstTransferDateTime}}
   PURPOSE: {{.FirstTransferPurpose}}
   CONDITION: {{.FirstTransferCondition}}
   NOTES: {{.FirstTransferNotes}}
   SIGNATURE (FROM): ___________________________
   SIGNATURE (TO): ___________________________

{{if .SecondCustodian}}
2. FROM: {{.FirstCustodian}}, {{.FirstCustodianTitle}}
   TO: {{.SecondCustodian}}, {{.SecondCustodianTitle}}
   DATE/TIME: {{.SecondTransferDateTime}}
   PURPOSE: {{.SecondTransferPurpose}}
   CONDITION: {{.SecondTransferCondition}}
   NOTES: {{.SecondTransferNotes}}
   SIGNATURE (FROM): ___________________________
   SIGNATURE (TO): ___________________________
{{end}}

FINAL DISPOSITION: {{.FinalDisposition}}
AUTHORIZED BY: {{.AuthorizedBy}}, {{.AuthorizedByTitle}}
DATE: {{.DispositionDate}}
`,
			TemplateVars: []string{
				"DepartmentName", "CaseNumber", "EvidenceNumber", "EvidenceDescription",
				"RecoveredBy", "RecoveredByBadge", "RecoveryLocation", "RecoveryDateTime",
				"RecoveryNotes", "FirstCustodian", "FirstCustodianTitle", "FirstTransferDateTime",
				"FirstTransferPurpose", "FirstTransferCondition", "FirstTransferNotes",
				"SecondCustodian", "SecondCustodianTitle", "SecondTransferDateTime",
				"SecondTransferPurpose", "SecondTransferCondition", "SecondTransferNotes",
				"FinalDisposition", "AuthorizedBy", "AuthorizedByTitle", "DispositionDate",
			},
			Department: "Evidence Unit",
			IsApproved: true,
		},
	}
}

// GetTemplateByName retrieves a template by name from the default templates
func GetTemplateByName(name string) *Template {
	templates := GetDefaultTemplates()
	for _, t := range templates {
		if t.Name == name {
			return t
		}
	}
	return nil
}

// GetTemplateByID retrieves a template by ID from the default templates
func GetTemplateByID(id string) *Template {
	templates := GetDefaultTemplates()
	for _, t := range templates {
		if t.ID == id {
			return t
		}
	}
	return nil
}

// GetTemplatesByType retrieves templates by type from the default templates
func GetTemplatesByType(correspondenceType CorrespondenceType) []*Template {
	templates := GetDefaultTemplates()
	var result []*Template
	for _, t := range templates {
		if t.Type == correspondenceType {
			result = append(result, t)
		}
	}
	return result
}

// GetTemplatesByDepartment retrieves templates by department from the default templates
func GetTemplatesByDepartment(department string) []*Template {
	templates := GetDefaultTemplates()
	var result []*Template
	for _, t := range templates {
		if t.Department == department || t.Department == "Any" {
			result = append(result, t)
		}
	}
	return result
}
