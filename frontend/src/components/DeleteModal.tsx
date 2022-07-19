import {
  useDisclosure,
  Button,
  AlertDialog,
  AlertDialogOverlay,
  AlertDialogContent,
  AlertDialogHeader,
  AlertDialogBody,
  AlertDialogFooter,
  IconButton,
} from "@chakra-ui/react";
import React, { FC } from "react";
import { MdDelete } from "react-icons/md";

export const DeleteModal: FC<{ onDelete: () => void }> = ({ onDelete }) => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const cancelRef = React.useRef() as any;

  return (
    <>
      <IconButton
        variant="ghost"
        aria-label="delete"
        icon={<MdDelete />}
        onClick={onOpen}
      />

      <AlertDialog
        isOpen={isOpen}
        leastDestructiveRef={cancelRef}
        onClose={onClose}
      >
        <AlertDialogOverlay>
          <AlertDialogContent>
            <AlertDialogHeader fontSize="lg" fontWeight="bold">
              Delete Story Pack
            </AlertDialogHeader>

            <AlertDialogBody>
              Are you sure you want to delete this pack? This action cannot be
              canceld (but you can reinstall the pack from this tool or the
              lunii software)
            </AlertDialogBody>

            <AlertDialogFooter>
              <Button ref={cancelRef} onClick={onClose}>
                Cancel
              </Button>
              <Button
                colorScheme="red"
                onClick={() => {
                  onDelete();

                  onClose();
                }}
                ml={3}
              >
                Delete
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialogOverlay>
      </AlertDialog>
    </>
  );
};
