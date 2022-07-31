import { useQuery } from "@tanstack/react-query";
import {
  CheckForUpdate,
  GetDeviceInfos,
  GetInfos,
} from "../../wailsjs/go/main/App";

export const useDeviceQuery = () =>
  useQuery(["device"], GetDeviceInfos, {
    staleTime: Infinity,
  });

export const useUpdateQuery = () =>
  useQuery(["update"], CheckForUpdate, {
    cacheTime: Infinity,
  });

export const useGetInfosQuery = () =>
  useQuery(["infos"], GetInfos, {
    cacheTime: Infinity,
  });
