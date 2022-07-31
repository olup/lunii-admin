import {
  Badge,
  Box,
  Button,
  Code,
  Divider,
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
import { FC, version } from "react";
import { useLocation, useRoute } from "wouter";
import { GetDeviceInfos } from "../../wailsjs/go/main/App";
import {
  useDeviceQuery,
  useGetInfosQuery,
  useUpdateQuery,
} from "../hooks/useApi";
import { formatBytes } from "../utils";

export const AboutModal: FC = () => {
  const { data: device } = useDeviceQuery();
  const { data: infos } = useGetInfosQuery();
  const { data: update } = useUpdateQuery();
  const [_, setLocation] = useLocation();

  return (
    <>
      <Modal isOpen={true} onClose={() => setLocation("/")}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>About</ModalHeader>
          <ModalBody>
            {device && (
              <>
                <Divider mb={3} />
                <Box mb={3}>
                  <Box fontWeight="bold" mb={2}>
                    Lunii
                  </Box>
                  <Box mb={2}>
                    Serial Number <Code>{device.serialNumber}</Code>
                  </Box>
                  <Box mb={2}>
                    Version{" "}
                    <Code>
                      {device.firmwareVersionMajor}.
                      {device.firmwareVersionMinor}
                    </Code>
                  </Box>
                  <Box mb={3}>
                    Space{" "}
                    <Code>
                      {formatBytes(device.diskUsage.used)} /{" "}
                      {formatBytes(device.diskUsage.free)}
                    </Code>
                  </Box>
                  <Progress
                    size="sm"
                    value={
                      (100 / device.diskUsage.free) * device.diskUsage.used
                    }
                  />
                </Box>
              </>
            )}
            <Divider mb={3} />
            <Box mb={3}>
              <Box fontWeight="bold" mb={2}>
                Software
              </Box>
              <Box mb={2}>
                Version <Code>{infos?.version}</Code>
              </Box>
              <Box mb={2}>
                Computer ID <Code>{infos?.machineId}</Code>
              </Box>
            </Box>
            {update?.canUpdate && (
              <>
                <Divider mb={3} />
                <Box>
                  <Badge colorScheme="green" mb={2}>
                    New version
                  </Badge>
                  <Box mb={2}>A new version is available</Box>
                  <Box mb={2} fontStyle="italic" opacity={0.5}>
                    Very soon you'll be able to trigger the update from here
                  </Box>
                </Box>
              </>
            )}
          </ModalBody>

          <ModalFooter></ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
};
