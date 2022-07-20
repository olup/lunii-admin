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

export const IsInstallingModal: FC<{ isOpen: boolean }> = ({ isOpen }) => {
  return (
    <>
      <Modal isOpen={isOpen} onClose={() => {}}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Installing Pack</ModalHeader>
          <ModalBody>
            <Text mb={2}>Wait a few moments, it's almost ready.</Text>
            <Progress isIndeterminate />
          </ModalBody>

          <ModalFooter></ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
};
