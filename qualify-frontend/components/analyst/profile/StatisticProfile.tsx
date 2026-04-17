import { IStatisticProfile } from '@/types/analyst/statisticProfile';

const StatisticProfileCardsVector: IStatisticProfile[] = [
    {
        value: "124",
        statistic: "Projetos Completos"
    },
    {
        value: "98%",
        statistic: "Taxa de sucesso"
    },
    {
        value: "R$120/hr",
        statistic: "Preço por hora"
    },
    {
        value: "1k+",
        statistic: "Bugs descobridos"
    },
];

function Card({ value, statistic }: IStatisticProfile) {
  return (
    <div className="flex justify-start inset-ring-2 inset-ring-zinc-800 flex-col w-full mx-3 mt-4 bg-gray-950">
        <p className="text-blue-800 text-2xl mt-4 mx-4">
            {value}
        </p>
        <p className="text-zinc-400 mt-1 mb-4 mx-4">
            {statistic}
        </p>
    </div>
  );
}

export function StatisticProfile(){
    return(
        <div className="flex flex-col justify-between content-center mt-8 w-2/10 h-100">
            {StatisticProfileCardsVector.map((StatisticProfileCard) => (
                <Card key={StatisticProfileCard.statistic} {...StatisticProfileCard} />
            ))}
        </div>
    );
}