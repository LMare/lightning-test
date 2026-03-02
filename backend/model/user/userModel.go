package user

type UserModel struct {
	ID     string `json:"id"`
	Nom    string `json:"nom"`
	Prenom string `json:"prenom"`
	Age    int    `json:"age"`
	Email  string `json:"email"`
}

func New(nom string, prenom string, age int) UserModel {
	return UserModel{Nom: nom, Prenom: prenom, Age: age}
}

func NewEmptyUserModel() *UserModel {
	return &UserModel{}
}

func (p *UserModel) SetNom(n string) *UserModel {
	p.Nom = n
	return p
}

func (p *UserModel) SetPrenom(n string) *UserModel {
	p.Prenom = n
	return p
}

func (p *UserModel) SetAge(a int) *UserModel {
	p.Age = a
	return p
}
