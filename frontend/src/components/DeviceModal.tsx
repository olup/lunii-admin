import {
  Box,
  Button,
  Code,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Progress,
  Text,
} from "@chakra-ui/react";
import { useQuery } from "@tanstack/react-query";
import { FC } from "react";
import { useLocation, useRoute } from "wouter";
import { GetDeviceInfos } from "../../wailsjs/go/main/App";
import { useDeviceQuery } from "../hooks/useApi";
import { formatBytes } from "../utils";

export const DeviceModal: FC = () => {
  const { data: device } = useDeviceQuery();
  const [_, setLocation] = useLocation();

  if (!device) return null;

  return (
    <>
      <Modal isOpen={true} onClose={() => setLocation("/")}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>My Lunii</ModalHeader>
          <ModalBody>
            <Box>
              <Box mb={2}>
                Serial Number <Code>{device.serialNumber}</Code>
              </Box>
              <Box mb={2}>
                Version{" "}
                <Code>
                  {device.firmwareVersionMajor}.{device.firmwareVersionMinor}
                </Code>
              </Box>
              <Box mb={2}>
                Space{" "}
                <Code>
                  {formatBytes(device.diskUsage.used)} /{" "}
                  {formatBytes(device.diskUsage.free)}
                </Code>
              </Box>
              <Progress
                size="sm"
                value={(100 / device.diskUsage.free) * device.diskUsage.used}
              />
            </Box>
          </ModalBody>

          <ModalFooter></ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
};
