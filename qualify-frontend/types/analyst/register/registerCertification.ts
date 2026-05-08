import { Certification } from "@/types/services/certification";

export interface IRegisterCertifications {
  certificationsAnalyst: Certification[];
  setCertificationsAnalyst: React.Dispatch<
    React.SetStateAction<Certification[]>
  >;
}
