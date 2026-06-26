"use client";

import { useState } from "react";
import { X } from "lucide-react";
import { FormInput, FormButton } from "@/components/ui";
import { proposalService } from "@/libs/services";

interface ProposalModalProps {
  analystId: number;
  clientId: number;
  onClose: () => void;
  onSuccess: () => void;
}

export function ProposalModal({
  analystId,
  clientId,
  onClose,
  onSuccess,
}: ProposalModalProps) {
  const [title, setTitle] = useState("");
  const [hourlyRate, setHourlyRate] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) {
      setError("Informe o serviço.");
      return;
    }
    const rate = parseFloat(hourlyRate.replace(",", "."));
    if (isNaN(rate) || rate <= 0) {
      setError("Informe um valor por hora válido.");
      return;
    }
    setError("");
    setLoading(true);
    const result = await proposalService.create({
      title: title.trim(),
      content: title.trim(),
      client_id: clientId,
      analyst_id: analystId,
      proposed_hourly_rate: rate,
    });
    setLoading(false);
    if (result) {
      onSuccess();
    } else {
      setError("Erro ao enviar proposta. Tente novamente.");
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
      <div className="w-full max-w-md rounded-2xl border border-white/10 bg-bg-dark p-8 shadow-2xl">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-bold text-white">Enviar Proposta</h2>
          <button
            onClick={onClose}
            className="text-neutral-slate hover:text-white transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <FormInput
            label="Serviço"
            name="title"
            placeholder="Ex: Sanity Test no Módulo X"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
          />
          <FormInput
            label="Valor por hora (R$)"
            name="hourly_rate"
            type="number"
            placeholder="150.00"
            min="0.01"
            step="0.01"
            value={hourlyRate}
            onChange={(e) => setHourlyRate(e.target.value)}
            required
          />

          {error && <p className="text-red-400 text-sm">{error}</p>}

          <div className="flex gap-3 mt-2">
            <FormButton
              type="button"
              variant="outline"
              onClick={onClose}
              disabled={loading}
            >
              Cancelar
            </FormButton>
            <FormButton type="submit" loading={loading} loadingText="Enviando...">
              Enviar Proposta
            </FormButton>
          </div>
        </form>
      </div>
    </div>
  );
}
