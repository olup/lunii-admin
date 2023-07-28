import {
  Button,
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
import { FC } from "react";

import { GetCurrentJob } from "../../wailsjs/go/main/App";
import { useQuery } from "@tanstack/react-query";

export const IsInstallingModal: FC<{ isOpen: boolean }> = ({ isOpen }) => {
  const { data } = useQuery(["currentJob"], () => GetCurrentJob(), {
    refetchInterval: 500,
  });

  if (!data) return null;

  return (
    <>
      <Modal isOpen={isOpen} onClose={() => {}}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Installing Pack</ModalHeader>
          <ModalBody>
            <Text>Unpacking</Text>
            {data.unpackDone && <Text>Converting images</Text>}
            {data.imagesConversionDone && <Text>Converting audios</Text>}
            {data.audiosConversionDone && <Text>Copying</Text>}
            {data.copyingDone && <Text>Finalizing</Text>}
            <ModalFooter></ModalFooter>
          </ModalBody>
        </ModalContent>
      </Modal>
    </>
  );
};
