import { Certification } from "@/types/services/certification";

export interface ICertificationSavesProps {
  certification: Certification;
  setCertificationsAnalyst: React.Dispatch<
    React.SetStateAction<Certification[]>
  >;
}
