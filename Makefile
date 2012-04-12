all: gedcom_main.6
	6l gedcom_main.6

gedcom.6: search/gedcom.go
	6g search/gedcom.go


gedcom_main.6: tools/gedcom_main.go gedcom.6
	6g tools/gedcom_main.go

