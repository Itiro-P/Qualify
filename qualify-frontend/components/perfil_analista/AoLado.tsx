interface tecnologiaCardsAnal {
    tecnologia: string;
}

interface certificacaoCardsAnal {
  nome: string;
  instituicao: string;
  ano: string;
}


const tecnologiasAnal: tecnologiaCardsAnal[] = [
    {
        tecnologia: "Docker",
    },
    {
        tecnologia: "Javascript/TS",
    },
    {
        tecnologia: "GitHub Actions",
    },
    {
        tecnologia: "Kubernets",
    },
];

const certificacoesAnal: certificacaoCardsAnal[] = [
  {
    nome: "Certificado de Testador Supremo de Botões",
    instituicao: "Irmandade do Teste Supremo",
    ano: "2024",
  },
  {
    nome: "Certificado de Guardião da Qualidade",
    instituicao: "Ordem Secreta do Teste",
    ano: "2023",
  },
  {
    nome: "Certificado Profissional em Testes Aleatórios",
    instituicao: "Comitê do Teste Aleatório",
    ano: "2023",
  }
];

function TecnologiasCards({ tecnologia }: tecnologiaCardsAnal) {
    return (
        <p className="text-blue-600 border-blue-600 rounded-xl text-lg px-2 py-1 m-1 bg-blue-950">{tecnologia}</p>
    );
}

function CertificacoesCards({ nome, instituicao, ano }: certificacaoCardsAnal) {
  return (
    <div className="flex flex-col mt-2 mb-4">
        <h3 className="text-white-950 mb-1">{nome}</h3>
        <p className="text-zinc-400 text-sm mb-1">{instituicao}</p>
        <span className="text-zinc-500 text-sm">{ano}</span>
    </div>
  );
}

export function AoLado(){
    return(
        <section id="aolado" className="flex flex-col px-6 md:px-20 py-24 w-5/20">
            <div id="tecnologias" className="flex flex-col border-zinc-800 border rounded-xl p-4 mb-6">
                <p className="text-zinc-400 text-xl my-2">Tecnologias</p>
                <div className="flex flex-wrap mt-2 mx-2">
                    {tecnologiasAnal.map((tecnologias) => (
                        <TecnologiasCards key={tecnologias.tecnologia} {...tecnologias} />
                    ))}
                </div>
            </div>
            <div id="certificações" className="flex flex-col border-zinc-800 border rounded-xl p-4">
                <p className="text-zinc-400 text-xl my-2">Certificações</p>
                <div className="flex flex-wrap mx-3">
                    {certificacoesAnal.map((certificados) => (
                        <CertificacoesCards key={certificados.nome} {...certificados} />
                    ))}
                </div>
            </div>
        </section>
    );
}