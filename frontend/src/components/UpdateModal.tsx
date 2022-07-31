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
import { Link, useLocation, useRoute } from "wouter";
import { CheckForUpdate, GetDeviceInfos } from "../../wailsjs/go/main/App";
import { useDeviceQuery, useUpdateQuery } from "../hooks/useApi";
import { formatBytes } from "../utils";

export const UpdateModal: FC = () => {
  const [location, setLocation] = useLocation();
  const { data } = useUpdateQuery();

  return (
    <>
      <Modal isOpen={location === "/update"} onClose={() => setLocation("/")}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Update available</ModalHeader>
          <ModalBody>
            <Box mb={2}>
              A new version of this tool is available.
              <Box>
                Very soon this will let you self-update. Until this time, you
                can download the latest version from the github repo.
              </Box>
            </Box>
            <Box mb={2}>
              Version <Code>{data?.latestVersion}</Code>
            </Box>
            {data?.releaseNotes && (
              <Code display="block" whiteSpace="pre" p={2}>
                {data.releaseNotes}
              </Code>
            )}
          </ModalBody>

          <ModalFooter></ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
};
