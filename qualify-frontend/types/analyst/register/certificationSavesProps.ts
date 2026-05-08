import { ICertification } from "@/types/analyst/profile/certification";

export interface ICertificationSavesProps {
  certification: ICertification;
  setCertificationsAnalyst: React.Dispatch<
    React.SetStateAction<ICertification[]>
  >;
}
