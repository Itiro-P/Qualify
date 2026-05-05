import { ITechnology } from "@/types/analyst/profile/technology";

export interface ITechnologySavesProps {
  technology: ITechnology;
  setTechnologiesAnalyst: React.Dispatch<React.SetStateAction<ITechnology[]>>;
}
