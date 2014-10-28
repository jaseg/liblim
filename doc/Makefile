
all: overview.pdf
.PHONY: all


overview.pdf: overview.tex overview.bib
	pdflatex overview.tex
	pdflatex overview.tex
	bibtex overview
	pdflatex overview.tex
	pdflatex overview.tex

