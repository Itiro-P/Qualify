import { api } from "@/libs/api";
import type {
  ProposalLetter,
  ProposalLetterResponse,
  ProposalLettersResponse,
  ProposalLetterUpdateRequest,
  ListProposalsParams,
} from "@/types/services/proposal";

export const proposalService = {
  list(params?: ListProposalsParams): Promise<ProposalLettersResponse> {
    const search = new URLSearchParams();
    if (params?.client_id) search.set("client_id", String(params.client_id));
    if (params?.analyst_id) search.set("analyst_id", String(params.analyst_id));
    const qs = search.toString();
    return api.get(`/proposals${qs ? `?${qs}` : ""}`);
  },

  getById(id: number): Promise<ProposalLetterResponse> {
    return api.get(`/proposals/${id}`);
  },

  create(data: ProposalLetter): Promise<ProposalLetterResponse> {
    return api.post("/proposals", data);
  },

  update(id: number, data: ProposalLetter): Promise<ProposalLetterResponse> {
    return api.put(`/proposals/${id}`, data);
  },

  patch(
    id: number,
    data: ProposalLetterUpdateRequest,
  ): Promise<ProposalLetterResponse> {
    return api.patch(`/proposals/${id}`, data);
  },

  delete(id: number): Promise<Record<string, string>> {
    return api.delete(`/proposals/${id}`);
  },
};
