package client

import (
	"github.com/CyCoreSystems/ari-proxy/v5/proxy"
	"github.com/CyCoreSystems/ari/v5"
)

type storedRecording struct {
	c *Client
}

func (s *storedRecording) List(filter *ari.Key) ([]*ari.Key, error) {
	return s.c.listRequest(&proxy.Request{
		Kind: "RecordingStoredList",
		Key:  filter,
	})
}

func (s *storedRecording) Get(key *ari.Key) *ari.StoredRecordingHandle {
	k, err := s.c.getRequest(&proxy.Request{
		Kind: "RecordingStoredGet",
		Key:  key,
	})
	if err != nil {
		s.c.log.Warn("failed to get stored recording for handle", "error", err)
		return ari.NewStoredRecordingHandle(key, s, nil)
	}
	return ari.NewStoredRecordingHandle(k, s, nil)
}

func (s *storedRecording) Data(key *ari.Key) (*ari.StoredRecordingData, error) {
	data, err := s.c.dataRequest(&proxy.Request{
		Kind: "RecordingStoredData",
		Key:  key,
	})
	if err != nil {
		return nil, err
	}
	return data.StoredRecording, nil
}

func (s *storedRecording) Copy(key *ari.Key, dest string) (*ari.StoredRecordingHandle, error) {
	h := ari.NewStoredRecordingHandle(key.New(ari.StoredRecordingKey, dest), s, nil)

	err := s.c.commandRequest(&proxy.Request{
		Kind: "RecordingStoredCopy",
		Key:  key,
		RecordingStoredCopy: &proxy.RecordingStoredCopy{
			Destination: dest,
		},
	})

	// NOTE: Always return the handle, even when we have an error
	return h, err
}

func (s *storedRecording) Delete(key *ari.Key) error {
	return s.c.commandRequest(&proxy.Request{
		Kind: "RecordingStoredDelete",
		Key:  key,
	})
}
