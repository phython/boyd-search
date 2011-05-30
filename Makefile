gedcom.6: gedcom.go
	/home/${USER}/hg/go/bin/6g gedcom.go


gedcom_main.6: gedcom_main.go gedcom.6
	/home/${USER}/hg/go/bin/6g gedcom_main.go

all: gedcom_main.6
	/home/${USER}/hg/go/bin/6l gedcom_main.6
