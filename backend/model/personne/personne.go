package personne

type Personne struct {
	Nom    string `json:"nom"`
	Prenom string `json:"prenom"`
	Age    int    `json:"age"`
}

func New(nom string, prenom string, age int) Personne {
	return Personne{Nom: nom, Prenom: prenom, Age: age}
}

func NewEmptyPersonne() *Personne {
	return &Personne{}
}

func (p *Personne) SetNom(n string) *Personne {
	p.Nom = n
	return p
}

func (p *Personne) SetPrenom(n string) *Personne {
	p.Prenom = n
	return p
}

func (p *Personne) SetAge(a int) *Personne {
	p.Age = a
	return p
}
