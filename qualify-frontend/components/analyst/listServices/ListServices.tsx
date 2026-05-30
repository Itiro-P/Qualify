"use client";

import { analystService } from "@/libs/services/analystService";
import { Analyst } from "@/types/services/analyst";
import { Service } from "@/types/services/service";
import { useEffect, useState } from "react";

function ListServicesArray({ services }: { services: Service[] }) {
  return (
    <div>
      {services.map((service) => (
        <div key={service.id}>
          <p>{service.title}</p>
          <p>{service.content}</p>
        </div>
      ))}
    </div>
  );
}

export function ListServices({ analyst }: { analyst: Analyst }) {
  const [servicesCompleteds, setServicesCompleteds] = useState<Service[]>([]);
  const [servicesProposal, setServicesProposal] = useState<Service[]>([]);
  const [servicesNegotiation, setServicesNegotiation] = useState<Service[]>([]);
  const [servicesBlocked, setServicesBlocked] = useState<Service[]>([]);
  const [servicesInProgress, setServicesInProgress] = useState<Service[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [abaAtiva, setAbaAtiva] = useState<
    "propostas" | "negociacao" | "em andamento" | "bloqueados" | "concluidos"
  >("concluidos");

  useEffect(() => {
    async function getInfo() {
      setLoading(true);
      const analystServices = await analystService.listServices(analyst.id);
      if (analystServices) {
        setServicesCompleteds(
          analystServices.filter((s) => s.status === "COMPLETED"),
        );
        setServicesProposal(
          analystServices.filter((s) => s.status === "PROPOSAL"),
        );
        setServicesNegotiation(
          analystServices.filter((s) => s.status === "NEGOTIATION"),
        );
        setServicesBlocked(
          analystServices.filter((s) => s.status === "BLOCKED"),
        );
        setServicesInProgress(
          analystServices.filter((s) => s.status === "IN_PROGRESS"),
        );
      } else {
        setError("Erro ao carregar os serviços do analista.");
      }
      setLoading(false);
    }
    getInfo();
  }, [analyst.id]);

  return (
    <div>
      <div className="flex">
        <button
          className="p-2 m-2 bg-blue-950"
          onClick={() => setAbaAtiva("propostas")}
        >
          Propostas
        </button>

        <button
          className="p-2 m-2 bg-blue-950"
          onClick={() => setAbaAtiva("negociacao")}
        >
          Negociação
        </button>

        <button
          className="p-2 m-2 bg-blue-950"
          onClick={() => setAbaAtiva("bloqueados")}
        >
          Bloqueados
        </button>

        <button
          className="p-2 m-2 bg-blue-950"
          onClick={() => setAbaAtiva("em andamento")}
        >
          Em Andamento
        </button>

        <button
          className="p-2 m-2 bg-blue-950"
          onClick={() => setAbaAtiva("concluidos")}
        >
          Concluídos
        </button>
      </div>

      <div className="mt-4">
        {abaAtiva === "propostas" && (
          <ListServicesArray services={servicesProposal} />
        )}

        {abaAtiva === "negociacao" && (
          <ListServicesArray services={servicesNegotiation} />
        )}

        {abaAtiva === "em andamento" && (
          <ListServicesArray services={servicesInProgress} />
        )}

        {abaAtiva === "bloqueados" && (
          <ListServicesArray services={servicesBlocked} />
        )}

        {abaAtiva === "concluidos" && (
          <ListServicesArray services={servicesCompleteds} />
        )}
      </div>
    </div>
  );
}
