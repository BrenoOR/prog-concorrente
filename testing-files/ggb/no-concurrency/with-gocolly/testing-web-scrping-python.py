import requests
from bs4 import BeautifulSoup

url = 'https://www.techspot.com/article/2547-what-are-threads/'
response = requests.get(url)

if response.status_code == 200:
    soup = BeautifulSoup(response.text, 'html.parser')
    # Exemplo: Extrair todos os textos das tags <p>
    paragraphs = soup.find_all('p')
    for p in paragraphs:
        print(p.text)
else:
    print('Erro ao obter a página:', response.status_code)
