package main

func main() {
	// Modelo proposto:
	// 1. Servidor:
	// 1.1. O servidor tem internamente uma cópia de todas as páginas que serão acessadas (Simulando o servidor real)
	// 1.2. O servidor recebe uma requisição de um cliente, e retorna a página solicitada (Simulando o servidor real)
	// 1.3. Serão 2 modelos de servidor: Um que utiliza protocolo TCP e outro que utiliza protocolo UDP
	// 1.4. Para se assemelhar às atividades anteriores, o servidor não irá fazer o scraping, isso será feito pelo cliente

	// 2. Cliente:
	// 2.1. O cliente deve ser o mais próximo possível do cliente da primeira atividade.
	// 2.2. O goquery só se utiliza do pacote net/html de go, então basta enviar o HTML para o cliente e fazer o scraping lá
	// 2.3. Proponho realizar testes tanto com o protocolo TCP quanto com o protocolo UDP,
	//visto que mesmo adaptado, o modelo funciona de forma diferente do modelo anterior
}
