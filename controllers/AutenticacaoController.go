package controllers

import (
	"net/http"
	"strconv"

	"github.com/tayron/agenda-contatos/models"

	"github.com/tayron/agenda-contatos/bootstrap/library/session"
	"github.com/tayron/agenda-contatos/bootstrap/library/template"
	"github.com/tayron/agenda-contatos/bootstrap/library/util"
)

// ValidarSessao - Verifica se tem usuário logado e redireciona para tela de login
func ValidarSessao(w http.ResponseWriter, r *http.Request) {
	usuario := session.GetDadoSessao("login", w, r)

	if usuario == "" {
		http.Redirect(w, r, "/login", 302)
	}
}

// Login - Exibe tela de login e faz loogin do usuário no banco
func Login(w http.ResponseWriter, r *http.Request) {
	flashMessage := template.FlashMessage{}

	usuario := session.GetDadoSessao("login", w, r)

	if usuario != "" {
		http.Redirect(w, r, "/", 302)
	}

	if r.Method == "POST" {
		r.ParseForm()
		login := r.PostForm.Get("login")
		senha := r.PostForm.Get("senha")
		senhaCriptografada, _ := util.GerarHashSenha(senha)

		usuarioModel := models.Usuario{
			Login: login,
			Senha: senhaCriptografada,
			Ativo: true,
		}

		usuario := usuarioModel.BuscarPorLoginStatus()

		if util.CompararSenhaComHash(senha, usuario.Senha) == true {
			idUsuario := strconv.Itoa(usuario.ID)
			session.SetDadoSessao("id_usuario", idUsuario, w, r)
			session.SetDadoSessao("nome_usuario", usuario.Login, w, r)
			session.SetDadoSessao("login", usuario.Login, w, r)

			http.Redirect(w, r, "/", 302)
		}

		flashMessage.Type, flashMessage.Message = template.ObterTipoMensagemAcessoNegado()
	}

	parametros := template.Parametro{
		System:       template.ObterInformacaoSistema(w, r),
		FlashMessage: flashMessage,
	}

	template.LoadView(w, "template/autenticacao/*.html", "loginPage", parametros)
}

// Logout - Limpa os dados da sessão e redireciona para a tela de login
func Logout(w http.ResponseWriter, r *http.Request) {
	session.ClearDadosSessao(w, r)
	http.Redirect(w, r, "/login", 302)
}
