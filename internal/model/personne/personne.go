package personne

type personne struct {
	nom string
	prenom string
	age int
}

func New(nom string, prenom string, age int) personne {
	return personne{nom : nom, prenom : prenom, age : age,}
}

func NewEmptyPersonne() *personne {
	return &personne{}
}

func (p *personne) SetNom(n string) *personne {
	p.nom = n
	return p
}

func (p *personne) SetPrenom(n string) *personne {
	p.prenom = n
	return p
}

func (p *personne) SetAge(a int) *personne {
	p.age = a
	return p
}