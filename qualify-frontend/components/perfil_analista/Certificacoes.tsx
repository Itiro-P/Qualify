import { GraduationCap } from 'lucide-react';

interface certificacaoCardsAnal {
    nome: string;
    instituicao: string;
    descricao: string;
    ano: string;
}

const certificacoesAnal: certificacaoCardsAnal[] = [
  {
    nome: "Certificado de Testador Supremo de Botões",
    instituicao: "Irmandade do Teste Supremo",
    descricao: "Me da certificações de testar botões",
    ano: "2024",
  },
  {
    nome: "Certificado de Guardião da Qualidade",
    instituicao: "Ordem Secreta do Teste",
    descricao: "Me da certificações de testar botões",
    ano: "2023",
  },
  {
    nome: "Certificado Profissional em Testes Aleatórios",
    descricao: "Me da certificações de testar botões",
    instituicao: "Comitê do Teste Aleatório",
    ano: "2023",
  }
];

function CertificacoesCards({ nome, instituicao, descricao, ano }: certificacaoCardsAnal) {
  return (
    <div className="flex flex-col mt-4 mb-4">
        <h3 className="text-white-950 mb-1">
            {nome}
        </h3>
        <p className="text-zinc-400 text-sm mb-1">
            {instituicao}
        </p>
        <p className="text-zinc-200 text-base mb-1">
            {descricao}
        </p>
        <span className="text-zinc-500 text-sm">
            {ano}
        </span>
    </div>
  );
}

export function Certificacoes(){
    return(
        <div id="certificações" className="flex flex-col p-4">
            <div className="flex">
                <GraduationCap className="text-blue-500 mr-2"/>
                <h2 className="text-xl">
                    Certificações
                </h2>
            </div>
            <div className="flex flex-col mx-3">
                {certificacoesAnal.map((certificados) => (
                    <CertificacoesCards key={certificados.nome} {...certificados} />
                ))}
            </div>
        </div>
    );
}