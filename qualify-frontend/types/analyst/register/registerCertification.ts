import { ICertification } from "@/types/analyst/profile/certification";

export interface IRegisterCertifications {
  certificationsAnalyst: ICertification[];
  setCertificationsAnalyst: React.Dispatch<
    React.SetStateAction<ICertification[]>
  >;
}
