import { useQuery } from "@tanstack/react-query";
import { GetDeviceInfos } from "../../wailsjs/go/main/App";

export const useDeviceQuery = () =>
  useQuery(["device"], GetDeviceInfos, {
    staleTime: Infinity,
  });
