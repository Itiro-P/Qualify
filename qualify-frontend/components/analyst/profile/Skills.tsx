import { Settings } from "lucide-react";
import { Skill } from "@/types/services";

function Card({ name }: Skill) {
  return (
    <p className="text-blue-600 border-blue-600 rounded-xl text-lg px-2 py-1 m-1 bg-blue-950">
      {name}
    </p>
  );
}

export function Skills({
  skillsCardsVector,
}: {
  skillsCardsVector: Skill[];
}) {
  return (
    <div className="flex flex-col p-4 mb-6">
      <div className="flex">
        <Settings className="text-blue-500 mr-2" />
        <h2 className="text-xl">Tecnologias</h2>
      </div>
      <div className="flex mt-2 mx-2 ">
        {skillsCardsVector.map((skills) => (
          <Card key={skills.id} {...skills} />
        ))}
      </div>
    </div>
  );
}
