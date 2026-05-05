import { Settings } from "lucide-react";
import { ITechnology } from "@/types/analyst/technology";

const technologiesCardsVector: ITechnology[] = [
  {
    technology: "Docker",
  },
  {
    technology: "Javascript/TS",
  },
  {
    technology: "GitHub Actions",
  },
  {
    technology: "Kubernets",
  },
];

function Card({ technology }: ITechnology) {
  return (
    <p className="text-blue-600 border-blue-600 rounded-xl text-lg px-2 py-1 m-1 bg-blue-950">
      {technology}
    </p>
  );
}

export function Technologies() {
  return (
    <div className="flex flex-col p-4 mb-6">
      <div className="flex">
        <Settings className="text-blue-500 mr-2" />
        <h2 className="text-xl">Tecnologias</h2>
      </div>
      <div className="flex mt-2 mx-2 ">
        {technologiesCardsVector.map((technologies) => (
          <Card key={technologies.technology} {...technologies} />
        ))}
      </div>
    </div>
  );
}
