import { TestBed } from '@angular/core/testing';

import { InterfacesService } from './interfaces.service';

describe('InterfacesService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: InterfacesService = TestBed.get(InterfacesService);
    expect(service).toBeTruthy();
  });
});
