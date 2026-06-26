import { api } from "@/libs/api";
import type {
  ProposalLetter,
  ProposalLetterCreateRequest,
  ProposalLetterResponse,
  ProposalLettersResponse,
  ProposalLetterUpdateRequest,
  ListProposalsParams,
} from "@/types/services/proposal";

export const proposalService = {
  list(params?: ListProposalsParams): Promise<ProposalLetter[] | null> {
    const search = new URLSearchParams();
    if (params?.client_id) search.set("client_id", String(params.client_id));
    if (params?.analyst_id) search.set("analyst_id", String(params.analyst_id));
    const qs = search.toString();
    return api
      .get<ProposalLettersResponse>(`/proposals${qs ? `?${qs}` : ""}`)
      .then(
        (resp) => {
          return resp.proposal_letters;
        },
        () => {
          return null;
        },
      );
  },

  getById(id: number): Promise<ProposalLetter | null> {
    return api.get<ProposalLetterResponse>(`/proposals/${id}`).then(
      (resp) => {
        return resp.proposal_letter;
      },
      () => {
        return null;
      },
    );
  },

  create(data: ProposalLetterCreateRequest): Promise<ProposalLetter | null> {
    return api.post<ProposalLetterResponse>("/proposals", data).then(
      (resp) => {
        return resp.proposal_letter;
      },
      () => {
        return null;
      },
    );
  },

  update(id: number, data: ProposalLetter): Promise<ProposalLetter | null> {
    return api.put<ProposalLetterResponse>(`/proposals/${id}`, data).then(
      (resp) => {
        return resp.proposal_letter;
      },
      () => {
        return null;
      },
    );
  },

  patch(
    id: number,
    data: ProposalLetterUpdateRequest,
  ): Promise<ProposalLetter | null> {
    return api.patch<ProposalLetterResponse>(`/proposals/${id}`, data).then(
      (resp) => {
        return resp.proposal_letter;
      },
      () => {
        return null;
      },
    );
  },

  delete(id: number): Promise<Record<string, string>> {
    return api.delete(`/proposals/${id}`);
  },

  listByClient(userId: number): Promise<ProposalLetter[] | null> {
    return api
      .get<ProposalLettersResponse>(`/users/${userId}/client/proposals`)
      .then(
        (resp) => resp.proposal_letters,
        () => null,
      );
  },

  listByAnalyst(userId: number): Promise<ProposalLetter[] | null> {
    return api
      .get<ProposalLettersResponse>(`/users/${userId}/analyst/proposals`)
      .then(
        (resp) => resp.proposal_letters,
        () => null,
      );
  },
};
