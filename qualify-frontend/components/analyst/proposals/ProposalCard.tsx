"use client";

import { useState, useEffect, useCallback } from "react";
import { FileText, Clock } from "lucide-react";
import { serviceService, proposalService } from "@/libs/services";
import type { ProposalLetter } from "@/types/services/proposal";
import type { Service } from "@/types/services/service";

interface ProposalCardProps {
  proposal: ProposalLetter;
  viewerRole: "client" | "analyst";
  onUpdate: () => void;
}

export function ProposalCard({
  proposal,
  viewerRole,
  onUpdate,
}: ProposalCardProps) {
  const [service, setService] = useState<Service | null>(null);
  const [loadingState, setLoadingState] = useState(true);
  const [acting, setActing] = useState(false);

  const fetchService = useCallback(async () => {
    if (!proposal.id) return;
    const services = await serviceService.list({ proposal_id: proposal.id });
    setService(services?.[0] ?? null);
    setLoadingState(false);
  }, [proposal.id]);

  useEffect(() => {
    fetchService();
  }, [fetchService]);

  const handleAccept = async () => {
    if (!proposal.id || !proposal.proposed_hourly_rate) return;
    setActing(true);
    await serviceService.create({
      proposal_letter_id: proposal.id,
      title: proposal.title ?? "",
      content: proposal.content ?? proposal.title ?? "",
      hourly_rate: proposal.proposed_hourly_rate,
      status: "Contrato Ativo",
    });
    setActing(false);
    onUpdate();
  };

  const handleReject = async () => {
    if (!proposal.id) return;
    setActing(true);
    await proposalService.delete(proposal.id);
    setActing(false);
    onUpdate();
  };

  const isPending = !loadingState && !service;

  const statusLabel = loadingState
    ? "Carregando..."
    : service
      ? "Contrato Ativo"
      : "Pendente";

  const statusClass = loadingState
    ? "bg-white/5 text-neutral-slate border-white/10"
    : service
      ? "bg-emerald-500/10 text-emerald-400 border-emerald-500/30"
      : "bg-yellow-500/10 text-yellow-400 border-yellow-500/30";

  return (
    <div className="rounded-xl border border-white/10 bg-white/5 p-6">
      <div className="flex items-start justify-between gap-4">
        <div className="flex items-center gap-2 min-w-0">
          <FileText className="w-5 h-5 text-accent shrink-0" />
          <h3 className="text-base font-bold text-white truncate">
            {proposal.title}
          </h3>
        </div>
        <span
          className={`shrink-0 rounded-full border px-3 py-1 text-xs font-medium ${statusClass}`}
        >
          {statusLabel}
        </span>
      </div>

      <p className="mt-3 flex items-center gap-1.5 text-sm text-neutral-slate">
        <Clock className="w-4 h-4 text-accent shrink-0" />
        R$ {proposal.proposed_hourly_rate?.toFixed(2)}/h
      </p>

      {proposal.time_created && (
        <p className="mt-1 text-xs text-neutral-slate">
          Enviada em{" "}
          {new Date(proposal.time_created).toLocaleDateString("pt-BR")}
        </p>
      )}

      {viewerRole === "analyst" && isPending && (
        <div className="flex gap-3 mt-5">
          <button
            onClick={handleAccept}
            disabled={acting}
            className="flex-1 rounded-lg bg-accent/10 border border-accent/30 text-accent px-4 py-2 text-sm font-medium hover:bg-accent/20 transition-colors disabled:opacity-50 cursor-pointer"
          >
            Aceitar
          </button>
          <button
            onClick={handleReject}
            disabled={acting}
            className="flex-1 rounded-lg bg-red-500/10 border border-red-500/30 text-red-400 px-4 py-2 text-sm font-medium hover:bg-red-500/20 transition-colors disabled:opacity-50 cursor-pointer"
          >
            Recusar
          </button>
        </div>
      )}
    </div>
  );
}
