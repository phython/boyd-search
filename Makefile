all: gedcom_main.6
	/home/${USER}/hg/go/bin/6l gedcom_main.6

gedcom.6: search/gedcom.go
	/home/${USER}/hg/go/bin/6g search/gedcom.go


gedcom_main.6: tools/gedcom_main.go gedcom.6
	/home/${USER}/hg/go/bin/6g tools/gedcom_main.go

