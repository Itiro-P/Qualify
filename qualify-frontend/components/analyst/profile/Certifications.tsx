import { GraduationCap } from 'lucide-react';
import { CertificationCard } from '@/types/analyst/certification';

const certificationsCardsVector: CertificationCard[] = [
  {
    name: "Certificado de Testador Supremo de Botões",
    institution: "Irmandade do Teste Supremo",
    description: "Me da certificações de testar botões",
    year: "2024",
  },
  {
    name: "Certificado de Guardião da Qualidade",
    institution: "Ordem Secreta do Teste",
    description: "Me da certificações de testar botões",
    year: "2023",
  },
  {
    name: "Certificado Profissional em Testes Aleatórios",
    description: "Me da certificações de testar botões",
    institution: "Comitê do Teste Aleatório",
    year: "2023",
  }
];

function Card({ name, institution, description, year }: CertificationCard) {
  return (
    <div className="flex flex-col mt-4 mb-4">
        <h3 className="text-white-950 mb-1">
            {name}
        </h3>
        <p className="text-zinc-400 text-sm mb-1">
            {institution}
        </p>
        <p className="text-zinc-200 text-base mb-1">
            {description}
        </p>
        <span className="text-zinc-500 text-sm">
            {year}
        </span>
    </div>
  );
}

export function Certifications(){
    return(
        <div id="certifications" className="flex flex-col p-4">
            <div className="flex">
                <GraduationCap className="text-blue-500 mr-2"/>
                <h2 className="text-xl">
                    Certificações
                </h2>
            </div>
            <div className="flex flex-col mx-3">
                {certificationsCardsVector.map((certification) => (
                    <Card key={certification.name} {...certification} />
                ))}
            </div>
        </div>
    );
}